Feature: Customer Management
  As a system administrator
  I want to manage customer information
  So that I can maintain customer records in the system

  Background:
    Given the customer service is running
    And the database is clean

  Scenario: Create a new customer
    When I send a request to create a customer with the following details:
      | name  | John Doe         |
      | email | john@example.com |
      | cpf   | 12345678901      |
    Then I should receive a response with status 201
    And the response should contain the customer ID
    And the response should contain customer details

  Scenario: Get customer by ID
    Given a customer exists with ID "507f1f77bcf86cd799439011"
    When I send a request to get customer with ID "507f1f77bcf86cd799439011"
    Then I should receive a response with status 200
    And the response should contain customer details

  Scenario: Get customer by non-existent ID
    When I send a request to get customer with ID "507f1f77bcf86cd799439999"
    Then I should receive a response with status 404
    And the response should contain an error message "Customer not found"

  Scenario: Get customer by CPF
    Given a customer exists with CPF "12345678901"
    When I send a request to get customer with CPF "12345678901"
    Then I should receive a response with status 200
    And the response should contain customer details

  Scenario: Get customer by non-existent CPF
    When I send a request to get customer with CPF "99999999999"
    Then I should receive a response with status 404
    And the response should contain an error message "Customer not found"

  Scenario: List all customers
    Given the following customers exist:
      | name     | email            | cpf         |
      | John Doe | john@example.com | 12345678901 |
      | Jane Doe | jane@example.com | 98765432100 |
    When I send a request to list all customers
    Then I should receive a response with status 200
    And the response should contain a list of 2 customers

  Scenario: Update customer information
    Given a customer exists with ID "507f1f77bcf86cd799439011"
    When I send a request to update customer with ID "507f1f77bcf86cd799439011" with the following details:
      | name  | John Updated         |
      | email | john_new@example.com |
    Then I should receive a response with status 200
    And the response should contain the updated customer details

  Scenario: Update non-existent customer
    When I send a request to update customer with ID "507f1f77bcf86cd799439999" with the following details:
      | name | John Updated |
    Then I should receive a response with status 404
    And the response should contain an error message "Customer not found"

  Scenario: Delete customer
    Given a customer exists with ID "507f1f77bcf86cd799439011"
    When I send a request to delete customer with ID "507f1f77bcf86cd799439011"
    Then I should receive a response with status 204

  Scenario: Delete non-existent customer
    When I send a request to delete customer with ID "507f1f77bcf86cd799439999"
    Then I should receive a response with status 404
    And the response should contain an error message "Customer not found"