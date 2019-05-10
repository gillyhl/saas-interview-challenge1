# ForgeRock SaaS Software Engineer Coding Challenge

This forked repo is my implementation of the ForgeRock SaaS Software Engineer Coding Challenge. In this solution, a user is able to upload a set of multiple-choice exam papers and get them marked (using a simple worker pattern).

The technical stack is as follows:

- Node.js / Express.js RESTful Web API to upload the exams to
- Redis for storage and pubsub
- Golang to listen for new exams posted and work on marking them.

## Building the solution

Once you have cloned the repo, run the following commands

```
docker-compose -f docker-compose.yml build
docker-compose -f docker-compose.yml up -d
```

All being well, all three services will be up and running.

## Usage instructions

When the solution is up, there will be a RESTful API running on port 5000, and the following HTTP requests can be made: All URLs are prefixed with `/api/v1/`
|URL|HTTP Verb|Returns|Comments|
|---|----|----|----|
|`/exams`|`POST`|Returns `id` of exam|Send a JSON file with the field key `exam` to start marking.|
|`/exams/:id`|`GET`|Gets all the marked papers for the given exam id.||
|`/exams/:id/papers/:id`|`GET`|Gets all the marked papers for the given exam id and paper id.||

Example cURL commands are (update the `-F` flag location to point to the exam file you want to process):

**NB**: These cURL commands were run through Window's command line. Using Postman is another viable option.

### Uploading exam

```curl
curl -X POST http://localhost:5000/api/v1/exams -F "exam=@C:\dev\saas\data\exam-1.json"
```

### Getting marked papers for exam

Update `<id>` to be the id returned from the above command

```curl
curl -X GET http://localhost:5000/api/v1/exams/<id>/
```

### Getting specific paper from exam

Update `<id>` to be the id returned from the above command

```curl
curl -X GET http://localhost:5000/api/v1/exams/<id>/papers/student-0
```

## Generating your own exam files

The `data` directory has a couple of example exam `JSON` files. You can generate your own by executing the following command:

```
node data/generateExam.js <examName> <numberOfQuestions> <numberOfPapers> <probability> <fileName>
```

| Argument          | Description                                                                                                                                 |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| examName          | String name of exam                                                                                                                         |
| numberOfQuestions | Number of questions to create for the exam                                                                                                  |
| numberOfPapers    | The number of people who sit the exam                                                                                                       |
| probability       | Probability that a question is answered correctly. Range 0-100. If this is set to 100, every question answered by everyone will be correct. |
| fileName          | Name of file to export the results to. File **MUST** be JSON format                                                                         |

## Testing

The code base has two testing suites. The first one is on the web server side, that uses mocha, with sinon as the mocking library and chai as the assertion library. Supertest helps to test API end points.

These tests can be run by navigating to the `web-server` directory and executing:

```
npm install
npm test
```

On the worker side, these tests are integration tests, and require you to run the following command from the project root:

```
docker-compose -f docker-compose-integration-test.yml up -d
```

And then navigating to the `worker` directory and executing:

```
go test ./... -count=1
```

## Shortcomings of Code

This code is very tightly coupled to the exam problem space. With a bit more refinement, I would create `Worker` and `Task` types that would have methods associated with them to carry out work. The `Worker` would have a handler function type within it to execute the work when a `Task` comes in, via a channel.

Due to my lack of experience with Golang, I opted out of doing the web server in Golang to save some time and decided to churn out a quick web server in Express.js

The worker portion of the code that makes use of Redis, wasn't coded particularly well when it comes to testability in isolation. If I was to redo this portion, I would create an interface for redis interactions that could then be mocked out in unit tests. As the code didn't lend itself well to unit testing, I tested the Golang module as an integration test, with redis running when being tested.

## Scaling The Solution

This solution could be scaled to have multiple instances of the Golang worker module on a cluster. Extra care will be needed when doing this to make sure worker tasks aren't duplicated when work comes in across the instances.

This could be achieved by putting the work to be done in a queue and then block popping the queue. `LPUSH` and `BRPOP` in the redis world can achieve this.

## Sequential vs Parallel Tasks

When it comes to executing tasks sequentially, the tasks can be simply iterated over using normal loops. Each task then is worked on one after another.

For parallel tasks, you would start a new process/thread (or goroutine in Golang) for each task. Golang's implementation of channels allows an easy way for goroutines to communicate with each other to keep a track of the tasks.
