Feature: Customer Authentication
  As a customer
  I want to authenticate with my CPF
  So that I can access the customer service system

  Background:
    Given the customer service is running

  Scenario: Successful customer authentication
    Given a customer exists with CPF "12345678901"
    When I send an authentication request with CPF "12345678901"
    Then I should receive a response with status 200
    And the response should contain a valid JWT token

  Scenario: Failed authentication with invalid CPF
    Given a customer exists with CPF "12345678901"
    When I send an authentication request with CPF "98765432100"
    Then I should receive a response with status 401
    And the response should contain an error message "Invalid credentials"

  Scenario: Failed authentication with missing CPF
    When I send an authentication request with CPF ""
    Then I should receive a response with status 400
    And the response should contain an error message "CPF is required"