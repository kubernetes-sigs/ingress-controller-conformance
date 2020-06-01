package cmd

import (
	"bytes"
	"os"

	"github.com/spf13/cobra"

	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/assets"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/k8s"

	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/util"
)

var applyIngressClass string

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

// addObjectAnnotations adds the given map of key-value pairs to
// the object as annotations. Any existing annotations for these keys
// will be overwritten.
func addObjectAnnotations(obj v1.Object, kv map[string]string) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	for k, v := range kv {
		annotations[k] = v
	}

	obj.SetAnnotations(annotations)
}

// complete verifies if ApplyOptions are valid and without conflicts by cribbing from (*apply.ApplyOptions)Complete().
func complete(o *apply.ApplyOptions, f util.Factory) error {
	var err error

	o.FieldManager = "ingress-controller-conformance"

	o.DryRunStrategy = util.DryRunNone
	o.ForceConflicts = false
	o.ServerSideApply = false

	// Don't bother recording the apply in an attribute.
	o.Recorder = &genericclioptions.NoopRecorder{}

	// allow for a success message operation to be specified at print time
	o.ToPrinter = func(operation string) (printers.ResourcePrinter, error) {
		o.PrintFlags.NamePrintFlags.Operation = operation
		if o.DryRunStrategy == util.DryRunClient {
			o.PrintFlags.Complete("%s (dry run)")
		}
		if o.DryRunStrategy == util.DryRunServer {
			o.PrintFlags.Complete("%s (server dry run)")
		}
		return o.PrintFlags.ToPrinter()
	}

	o.DynamicClient, err = f.DynamicClient()
	if err != nil {
		return err
	}

	// NOTE(jpeach): Here we omit calling (*FilenameOptions)
	// RequireFilenameOrKustomize() which enforces that either the
	// '-f' or '-k' flags must be present. We manually specified
	// the resources to parse from the compiled assets, so we don't
	// need either.

	o.DeleteOptions = o.DeleteFlags.ToOptions(o.DynamicClient, o.IOStreams)

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
		applyOpts := apply.NewApplyOptions(newStandardIO())
		getter := k8s.NewClientGetter()

		if err := complete(applyOpts, util.NewFactory(getter)); err != nil {
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
		for _, name := range assets.AssetNames() {
			applyOpts.Builder = applyOpts.Builder.Stream(
				bytes.NewBuffer(assets.MustAsset(name)), name)
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
					addObjectAnnotations(obj, map[string]string{
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
	applyCmd.Flags().StringVar(&applyIngressClass,
		"ingress-class", "",
		"Ingress class to set on Ingress resources")

	rootCmd.AddCommand(applyCmd)
}
