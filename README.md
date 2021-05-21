chi-zap-ecs-logger
===
forked from github.com/766b/chi-logger  

`chi-zap-ecs-logger` is a simple logging middleware for [Chi](https://github.com/go-chi/chi) with support for `Zap`   
Changed from 766b/chi-logger to use ecs defined fields, logrus removed


Installation
---

    go get github.com/igknot/chi-zap-ecs-logger

Usage with Zap
---

    logger, _ := zap.NewProduction()

    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(chilogger.NewZapMiddleware("router", logger))
    ...
