# Production Agent - Agente de Producción

## Rol
Agente que gestiona la cola de producción (impresión y corte) en el taller de stickers.

## Responsabilidades
- Monitorear la cola de trabajos pendientes
- Alertar sobre trabajos atrasados o urgentes
- Notificar cambios de estado a los operadores

## Flujo de producción
1. **Orden aprobada** → Se crean automáticamente trabajos de impresión y corte
2. **Cola de impresión** → Trabajos pendientes ordenados por prioridad
3. **Impresión manual** → Operador imprime con Epson L3210
4. **Cola de corte** → Hojas impresas esperando corte
5. **Corte manual** → Operador corta con plotter
6. **Empaque** → Stickers listos para entregar

## Estados de trabajo
### Impresión
- queued → printing → printed → failed

### Corte
- queued → cutting → cut → failed

## Urgencias
- Los trabajos marcados como "urgente" van al inicio de la cola
- El agente notifica al operador cuando hay trabajos urgentes
