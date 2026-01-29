# go-code-snippet

Ce dépôt regroupe plusieurs exemples et patterns Go autour de la concurrence, des goroutines, du pooling, des higher-order functions, etc.

## Structure des dossiers

```
.
├── goroutine
│   ├── fanout-fanin
│   │   ├── main.go
│   │   └── readme.md
│   ├── worker
│   │   ├── main.go
│   │   └── schema.png
│   └── worker-pool
│       ├── main.go
│       ├── main.go.errgroup
│       ├── readme.md
│       └── schema-worker-pool.png
├── higher-order-functions
│   ├── decorator-example
│   │   ├── main.go
│   │   └── schema.png
│   ├── middleware-example
│   │   ├── main.go
│   │   └── schema.png
│   ├── middleware-example-2
│   │   ├── main.go
│   │   └── schema.png
│   ├── pipeline-system
│   │   ├── main.go
│   │   └── schema.png
│   └── type-constraint
│       └── main.go
├── pooling
│   ├── pool-example
│   │   └── main.go
│   ├── pool-example-benchmark
│   │   └── main.go
│   └── readme.md
└── tests
    ├── data-processing-pipeline
    │   ├── main.go
    │   └── go.mod
    └── middleware-retrial
        ├── error.go
        ├── go.mod
        ├── implementation
        │   ├── storage
        │   │   ├── dynamodb
        │   │   │   └── dynamodb.go
        │   │   └── mysql
        │   │       └── mysql.go
        ├── infrastructure
        │   └── storage
        │       ├── auth.go
        │       └── storage.go
        └── main.go
```

Chaque dossier contient un exemple autonome, souvent accompagné d'un schéma explicatif et d'un README local.

## Objectif

- Illustrer des patterns Go modernes et idiomatiques
- Fournir des bases pour des architectures concurrentes robustes
- Servir de support pédagogique ou de base pour de nouveaux projets

---

Pour plus de détails, consulte chaque dossier et ses fichiers README associés.
