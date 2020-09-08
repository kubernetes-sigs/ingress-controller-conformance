const report = require("./cucumber-html-reporter");
const assert = require("assert");

assert(process.env.INPUT_DIRECTORY, "Environment variable INPUT_DIRECTORY is not optional");
assert(process.env.OUTPUT_DIRECTORY, "Environment variable OUTPUT_DIRECTORY is not optional");

report.generate({
  jsonDir: process.env.INPUT_DIRECTORY,
  reportPath: process.env.OUTPUT_DIRECTORY,
  pageFooter: '<p><a href="https://github.com/kubernetes-sigs/ingress-controller-conformance">Kubernetes ingress controller conformance</a></p>',
  ingress: {
    controller: process.env.INGRESS_CONTROLLER || 'N/A',
    version: process.env.CONTROLLER_VERSION || 'N/A'
  }
});
