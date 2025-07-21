# Technology Context

## 1. Technologies Used

### Version 0 (Mockup)
- **Frontend:** HTML, JavaScript, Bootstrap, CSS
- **Backend:** Go
- **Database:** SQLite

### Version 1 (Production)
- **Database:** PostgreSQL
- **Deployment:** Azure App Service (as a monorepo for frontend and backend)

## 2. Development Setup
- For v0, the application will be developed and run on `localhost`.
- The project will eventually be structured as a monorepo to contain both the frontend and backend code.

## 3. Technical Constraints & Decisions
The initial `projectBrief.md` raised several technical questions. Here are the initial recommendations for v0:

- **Go Framework:** For the v0 mockup, using Go's built-in `net/http` package is sufficient and recommended. It keeps the application lightweight and free of external dependencies. For v1, a lightweight framework like **Gin** or **Echo** could be adopted if routing and middleware needs become more complex.
- **Caching (Redis):** Caching with Redis is not necessary for the v0 mockup. This is a consideration for v1 to improve performance once the application is handling significant traffic and repeated queries.
- **Async Processing (RabbitMQ/Kafka):** For v0, message queuing systems are overkill. Go's native concurrency features (**goroutines** and **channels**) are perfectly suited to handle the asynchronous calls to the SLIK OJK API. This approach is simpler to implement and manage for the initial version. We can re-evaluate the need for a dedicated message broker like RabbitMQ for v1 if the async workload becomes very complex.
