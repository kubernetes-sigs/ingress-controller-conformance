@sig-network @conformance @release-1.19
Feature: Path rules

  An Ingress may define routing rules based on the request path.

  If the HTTP request path matches one of the paths in the
  Ingress objects, the traffic is routed to its backend service.

  Background:
    Given an Ingress resource in a new random namespace
      """
      apiVersion: networking.k8s.io/v1
      kind: Ingress
      metadata:
        name: path-rules
      spec:
        rules:
          - host: "path-rules"
            http:
              paths:
                - path: /foo
                  pathType: Prefix
                  backend:
                    service:
                      name: path-rules-foo
                      port:
                        number: 80

                - path: /foo
                  pathType: Exact
                  backend:
                    service:
                      name: path-rules-exact
                      port:
                        number: 80
      """
    Then The Ingress status shows the IP address or FQDN where it is exposed

  Scenario: An Ingress with a prefix path rule without a trailing slash should send traffic to the matching backend service
  (exact /foo matches request /foo)
    When I send a "GET" request to "http://path-rules/foo"
    Then the response status-code must be 200
    And the response must be served by the "path-rules-exact" service
    And the request path must be "/foo"

  Scenario: An Ingress with a prefix path rule with a trailing slash should send traffic to the matching backend service
  (prefix /foo matches request /foo/)
    When I send a "GET" request to "http://path-rules/foo/"
    Then the response status-code must be 200
    And the response must be served by the "path-rules-foo" service
    And the request path must be "/foo/"