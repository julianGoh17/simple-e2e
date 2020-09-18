# Simple-E2E

This project seeks to solve the problem of "How do we make testing easier and modular?". End to End testing is a useful tool that gives developers the confidence that their code is functioning and prevents regressions. So why is it so hard to write understandable tests and maintain the codebase? This project aims to solve these problems and make testing easy to understand.

## Components of Simple-E2E

A **test step** is the basic unit of a test and it consists of a string that clearly describes what you are trying to do to and a function that contains the operational code that it will perform. **Test steps** can be organized into **stages** which describe a set of **test steps** that aim to do similar things. For example, you could have a test step to initialize an API and another test step that adds data into a database in a  'set up' stage. Various stages can then combined to create a clear and understandable test. See [this example test](https://github.com/julianGoh17/simple-e2e/blob/master/tests/examples/multi-stage-test.yaml) to see how easy it is to understand a test in Simple-E2E!

Many more features are currently in the process of being planned and developed. To see what is being currently being planned and worked on, have a look at our [Kanban board](https://github.com/julianGoh17/simple-e2e/projects/1)!

## Prerequisites

If you are interested in using the project, the following are required to run the project:

- **Docker:** Allows us to run the project on any machine! You can download the necessary binaries from [the Docker Getting Started webpage](https://www.docker.com/get-started)
- **Golang**: The project is written in Golang and to utilize the framework you will need to write Golang code. The binary and some useful tutorials can be found on the [Golang website](https://golang.org/)

## Want to help out?

Any help is appreciated and welcome! The best way to get started is to get in touch with me or to have a look at our [Github issues](https://github.com/julianGoh17/simple-e2e/issues)! Have a read of the Github issues, pick one that interests you, and create a pull request!

## Credit where credit is due

Thank you to the team at IBM Event Streams who originally inspired me to come up with this idea. I could not have done it without you guys!
