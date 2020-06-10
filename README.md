# SLR Automation Backend

This repository contains the API which is used to manage Systematic Literature Review(SLR) projects.

# Running locally

Assumming you've installed [Golang](https://golang.org/doc/install)

1. Copy the [env example](src/slr-api/example.env) and setup your own environment
2. Make sure the environment is available in your terminal(`source yourfile.env`)
3. Execute: `make run`

# Docker

1. Build using `make docker-build`
2. Running the image: `docker run -e JWT_KEY=$JWT_KEY -e JWT_KEY_PUB=$JWT_KEY_PUB -e PG_HOST=$PG_HOST -e PG_PASSWORD=$PG_PASSWORD -e PG_USERNAME=$PG_USERNAME -p <your-exposing-port>:8000 slr-api` 

# Contributions

Got any improvements? Feel free to create a Pull request.

# Issues

Any issues? Please create an [Issue](https://github.com/lit-automation/backend/issues), or contact me personally at w.j.spaargaren@student.tudelft.nl.

# Related projects

The SLR automation environment consists of three individual tools, including this project.

* [Frontend project](https://github.com/lit-automation/frontend)
* [Chrome plugin](https://github.com/lit-automation/chrome-plugin)


# License

Licensed under the [MIT](LICENSE) license.

Created by [Wim Spaargaren](https://github.com/wimspaargaren)