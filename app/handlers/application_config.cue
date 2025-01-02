TemporalHostPort: "localhost:7233"
DBHost:       "localhost"
DBPort:       "5432"
DBUser:       "billing_service"
DBPassword:   "billing_service"
DBName:       "billing_service"
DBSchemaMigrationsPath: "db/migrations"


if #Meta.Environment.Name == "test" {
    // On this environment, we want to set ReadOnlyMode to true
    DBPort:       "5434"
    DBUser:       "billing_service_test"
    DBPassword:   "billing_service_test"
    DBName:       "billing_service_test"
    DBSchemaMigrationsPath: "migrations"
}