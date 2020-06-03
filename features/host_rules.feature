@sig-network @conformance @release-1.19
Feature: Host rules

  An Ingress may define routing rules based on the request host.

  If the HTTP request host matches one of the hosts in the
  Ingress objects, the traffic is routed to its backend service.

  Background:
    Given a new random namespace
    Given a self-signed TLS secret named "conformance-tls" for the "foo.bar.com" hostname
    Given an Ingress resource
      """
      apiVersion: networking.k8s.io/v1
      kind: Ingress
      metadata:
        name: host-rules
      spec:
        tls:
          - hosts:
              - foo.bar.com
            secretName: conformance-tls
        rules:
          - host: "*.bar.com"
            http:
              paths:
                - path: /
                  backend:
                    service:
                      name: host-rules-wildcard-bar-com
                      port:
                        number: 80
          - host: foo.bar.com
            http:
              paths:
                - backend:
                    service:
                      name: host-rules-foo-bar-com
                      port:
                        name: http
      """
    Then The Ingress status shows the IP address or FQDN where it is exposed

  Scenario: An Ingress with a host rule should send traffic to the matching backend service
    When I send a "GET" request to "https://foo.bar.com"
    Then the secure connection must verify the "foo.bar.com" hostname
    And the response status-code must be 200
    And the response must be served by the "host-rules-foo-bar-com" service
    And the request host must be "foo.bar.com"

  Scenario: An Ingress with a wildcard host rule should send traffic to the matching backend service
    When I send a "GET" request to "http://wildcard.bar.com"
    Then the response status-code must be 200
    And the response must be served by the "host-rules-wildcard-bar-com" service
    And the request host must be "wildcard.bar.com"