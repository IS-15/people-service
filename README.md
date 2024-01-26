# people-service

1. REST API Service to handle people (person)
2. GET with pagination and filtering
3. 
    a. POST to save new person.
    b. While handling POST gets additional data from age service, gender service, nationality(country) service.
    c. Has the simplest mock for the external services above for the case there is no access to the service.
4. DELETE person by id
5. PUT for person update.
6. Uses PostgreSQL.
7. Has migrations and simple migrator.
8. .env parameters are loaded in main() as it is a study case.