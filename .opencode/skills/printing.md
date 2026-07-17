---
name: printing
description: "Gestión de impresión Epson L3210 y plotter de corte"
---

## Módulo de Impresión

### Filosofía
La impresión y el corte son **100% manuales**. El sistema solo:
1. Muestra la cola de trabajos pendientes
2. Permite al admin marcar estados
3. Notifica al cliente cuando está listo

### Archivos
- `backend/internal/domain/print_job.go` - Modelo PrintJob
- `backend/internal/application/printing_service.go` - Lógica de cola
- `backend/internal/infra/repository/database/print_job_repository.go` - Persistencia
- `backend/internal/infra/entrypoint/rest/printing_controller.go` - Endpoints

### Tabla
- `print_jobs` - Cola de trabajos

### Estados de PrintJob
- print: queued → printing → printed → failed
- cut: queued → cutting → cut → failed

### Flujo de trabajo del admin
1. Abre panel de producción en el frontend
2. Ve cola de impresión (pendientes)
3. Abre el PDF, lo imprime con Epson L3210
4. Marca "Impreso" en el sistema
5. Lleva las hojas al plotter
6. Abre el mismo PDF, corta con plotter
7. Marca "Cortado" en el sistema
8. Empaqueta y marca "Listo"

### Hardware
- **Impresora**: Epson L3210 (tinta continua, A4)
- **Plotter**: Corte por HPGL/serial (marca según setup)
- Ver `hardware/printer/SETUP.md` y `hardware/plotter/SETUP.md`
