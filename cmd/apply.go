package cmd

import (
	"bytes"
	"fmt"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/apiversion"
	"os"

	"github.com/spf13/cobra"

	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/assets"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/k8s"

	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/delete"
	"k8s.io/kubectl/pkg/cmd/util"
)

var (
	applyIngressAPIVersion string
	applyIngressClass      string
	applyIngressController string
)

func newStandardIO() genericclioptions.IOStreams {
	return genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

// addObjectLabels adds the given map of key-value pairs to
// the object as labels. Any existing labels for these keys
// will be overwritten.
func addObjectLabels(obj v1.Object, kv map[string]string) {
	labels := obj.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}

	for k, v := range kv {
		labels[k] = v
	}

	obj.SetLabels(labels)
}

// replaceObjectAnnotations sets the given map of key-value pairs to
// the object as annotations if the annotation key is defined in the
// object.
func replaceObjectAnnotations(obj v1.Object, kv map[string]string) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	for k, v := range kv {
		if annotations[k] != "" {
			annotations[k] = v
		}
	}

	obj.SetAnnotations(annotations)
}

// completeApply verifies if ApplyOptions are valid and without conflicts by cribbing from (*apply.ApplyOptions)Complete().
func completeApply(o *apply.ApplyOptions, f util.Factory) error {
	var err error

	o.FieldManager = "ingress-controller-conformance"

	o.DryRun = false
	o.ForceConflicts = false
	o.ServerDryRun = false
	o.ServerSideApply = false

	// Don't bother recording the apply in an attribute.
	o.Recorder = &genericclioptions.NoopRecorder{}

	// allow for a success message operation to be specified at print time
	o.ToPrinter = func(operation string) (printers.ResourcePrinter, error) {
		o.PrintFlags.NamePrintFlags.Operation = operation
		if o.DryRun {
			o.PrintFlags.Complete("%s (dry run)")
		}
		if o.ServerDryRun {
			o.PrintFlags.Complete("%s (server dry run)")
		}
		return o.PrintFlags.ToPrinter()
	}

	o.DiscoveryClient, err = f.ToDiscoveryClient()
	if err != nil {
		return err
	}

	dynamicClient, err := f.DynamicClient()
	if err != nil {
		return err
	}

	// NOTE(jpeach): Here we omit calling (*FilenameOptions)
	// RequireFilenameOrKustomize() which enforces that either the
	// '-f' or '-k' flags must be present. We manually specified
	// the resources to parse from the compiled assets, so we don't
	// need either.

	o.DeleteOptions = o.DeleteFlags.ToOptions(dynamicClient, o.IOStreams)

	o.OpenAPISchema, _ = f.OpenAPISchema()
	o.Validator, err = f.Validator(true)
	if err != nil {
		return err
	}

	o.Builder = f.NewBuilder()
	o.Mapper, err = f.ToRESTMapper()
	if err != nil {
		return err
	}

	o.DynamicClient, err = f.DynamicClient()
	if err != nil {
		return err
	}

	o.Namespace, o.EnforceNamespace, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	return nil
}

// completeDelete verifies if DeleteOptions are valid and without conflicts by cribbing from (*delete.DeleteOptions)Complete().
func completeDelete(o *delete.DeleteOptions, f util.Factory, kinds string) error {
	r := f.NewBuilder().
		Unstructured().
		ContinueOnError().
		LabelSelectorParam(o.LabelSelector).
		FieldSelectorParam(o.FieldSelector).
		SelectAllParam(o.DeleteAll).
		AllNamespaces(o.DeleteAllNamespaces).
		ResourceTypeOrNameArgs(false, []string{kinds}...).RequireObject(false).
		Flatten().
		Do()
	err := r.Err()
	if err != nil {
		return err
	}
	o.Result = r

	dynamicClient, err := f.DynamicClient()
	if err != nil {
		return err
	}
	o.DynamicClient = dynamicClient

	o.Mapper, err = f.ToRESTMapper()
	if err != nil {
		return err
	}

	return nil
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Ingress conformance resources to the current cluster",
	Long: `Apply Ingress conformance resources to the current cluster.

This command applies Kubernetes resources that are needed for
conformance verifications to the currently active Kubernetes
cluster.  If the '--ingress-class' is given, then all Ingress
resources are annotated with the specified Ingress class.

Resources created by this command are labeled with the following
information:

    app.kubernetes.io/part-of=ingress-controller-conformance
    app.kubernetes.io/managed-by=ingress-controller-conformance
    app.kubernetes.io/version=$VERSION
`,

	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		getter := k8s.NewClientGetter()
		factory := util.NewFactory(getter)

		fmt.Print("cleaning managed resources from previous run... ")

		deleteOpts := &delete.DeleteOptions{
			IOStreams:           newStandardIO(),
			LabelSelector:       "app.kubernetes.io/managed-by=ingress-controller-conformance",
			DeleteAllNamespaces: true,
			Cascade:             true,
		}
		// Attempt to delete managed resources
		if err := completeDelete(deleteOpts, factory, "deployments,services,ingresses,secrets"); err != nil {
			return err
		}
		if err := deleteOpts.RunDelete(factory); err != nil {
			return err
		}

		// Attempt to delete managed ingresclasses, but don't fail if this cluster does not have this resource type
		if err := completeDelete(deleteOpts, factory, "ingressclasses"); err != nil {
			fmt.Println(err)
		} else if err := deleteOpts.RunDelete(factory); err != nil {
			fmt.Println(err)
		}

		applyOpts := apply.NewApplyOptions(newStandardIO())
		if err := completeApply(applyOpts, util.NewFactory(getter)); err != nil {
			return err
		}

		// Prefab the Builder so that when we add the asset
		// streams, they will pick up the right state.
		applyOpts.Builder = applyOpts.Builder.
			Unstructured().
			Schema(applyOpts.Validator).
			ContinueOnError().
			NamespaceParam(applyOpts.Namespace).
			DefaultNamespace()

		// Add all the packaged assets as streams to the
		// apply Builder.
		dir := fmt.Sprintf("deployments/%s", applyIngressAPIVersion)
		dirAssets, err := assets.AssetDir(dir)
		if err != nil {
			return fmt.Errorf("no assets found for apiVersion %v", applyIngressAPIVersion)
		}
		fmt.Printf("applying assets from %s %v\n", dir, dirAssets)

		if applyIngressController != "" {
			ingressClass := fmt.Sprintf(`apiVersion: %v
kind: IngressClass
metadata:
  name: conformance
  annotations:
    ingressclass.kubernetes.io/is-default-class: "true"
spec:
  controller: %v`, applyIngressAPIVersion, applyIngressController)
			applyOpts.Builder = applyOpts.Builder.Stream(bytes.NewBuffer([]byte(ingressClass)), "IngressClass")
		}

		for _, name := range dirAssets {
			buffer := bytes.NewBuffer(assets.MustAsset(dir + "/" + name))
			applyOpts.Builder = applyOpts.Builder.Stream(buffer, name)
		}

		// Run the builder.
		res := applyOpts.Builder.Flatten().Do()
		if err := res.Err(); err != nil {
			return err
		}

		infos, err := res.Infos()
		if err != nil {
			return err
		}

		for _, info := range infos {
			obj, err := meta.Accessor(info.Object)
			if err != nil {
				return err
			}

			addObjectLabels(obj, map[string]string{
				"app.kubernetes.io/part-of":    "ingress-controller-conformance",
				"app.kubernetes.io/managed-by": "ingress-controller-conformance",
				"app.kubernetes.io/version":    VERSION,
			})

			if applyIngressClass != "" {
				objType, err := meta.TypeAccessor(info.Object)
				if err != nil {
					return err
				}

				if objType.GetKind() == "Ingress" {
					replaceObjectAnnotations(obj, map[string]string{
						"kubernetes.io/ingress.class": applyIngressClass,
					})
				}
			}
		}

		applyOpts.SetObjects(infos)

		return applyOpts.Run()
	},
}

func init() {
	applyCmd.Flags().StringVar(&applyIngressAPIVersion,
		"api-version", "",
		fmt.Sprintf("apiVersion of resources to apply %s", apiversion.All))

	applyCmd.Flags().StringVar(&applyIngressClass,
		"ingress-class", "",
		"Ingress class to set on Ingress resources")

	applyCmd.Flags().StringVar(&applyIngressController,
		"ingress-controller", "",
		"inject an IngressClass resource with for a given spec.controller value")

	if err := applyCmd.MarkFlagRequired("api-version"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(applyCmd)
}
