# Sprint Retrospective – VectraG

## ¿Qué salió bien?
- La arquitectura modular (clean architecture) permitió desarrollar el CLI y el panel admin de forma independiente, facilitando las pruebas.
- El uso de Cobra y promptui agilizó la creación de los comandos interactivos.
- Dockerizar la aplicación con multi‑stage build resultó en una imagen final pequeña (alpine) y reproducible.
- El pipeline CI/CD con GitHub Actions detectó fallos de compilación tempranos.

## ¿Qué salió mal?
- No se planificaron tests de infraestructura (HTTP handlers, base de datos), lo que dejó esos paquetes sin cobertura.
- La integración del frontend React empaquetado dentro del binario Go requirió ajustes inesperados en el Dockerfile.
- La implementación de S3 quedó como infraestructura preparada pero sin integrar aún en el código, lo cual retrasa la funcionalidad completa de gestión de medios.

## ¿Qué haremos diferente?
- Incluir desde el primer sprint la configuración de infraestructura como código (Terraform) para alinear desarrollo y operaciones.
- Escribir tests de integración para endpoints HTTP usando `httptest`.
- Planificar mejor las dependencias entre frontend y backend antes del build multi‑stage.
