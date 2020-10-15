@sig-network @conformance @release-1.19
Feature: Load Balancing
  An Ingress exposing a backend service with multiple replicas should use all the pods available
  The feature sessionAffinity is not configured in the backend service https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#service-v1-core

  Background:
    Given a new random namespace
    Given an Ingress resource named "load-balancing" with this spec:
    """
    defaultBackend:
      service:
        name: echo-service
        port:
          number: 8080
    """
    Then The Ingress status shows the IP address or FQDN where it is exposed
    Then The backend deployment "echo-service" for the ingress resource is scaled to 10

  Scenario Outline: An Ingress with no rules should send all requests to the default backend and
    When I send 100 requests to "http://load-balancing"
    Then all the responses status-code must be 200 and the response body should contain the IP address of 10 different Kubernetes pods
