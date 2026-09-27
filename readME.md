crear migraciones:
```
migrate create -ext sql -dir migrations -seq add_nombre_migracion_table

# Para aplicar las migraciones (UP)
migrate -path migrations -database "postgres://postgres:admin@localhost:5432/content_db?sslmode=disable" up

# Para revertir la última migración (DOWN)
migrate -path migrations -database "postgres://postgres:admin@localhost:5432/content_db?sslmode=disable" down 1
```