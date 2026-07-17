# Configuración del Plotter de Corte

## Conexión
- El plotter se conecta por USB o puerto serial
- Protocolo común: HPGL (Hewlett-Packard Graphics Language)
- Velocidad de corte: ajustable según material

## Drivers
- Depende de la marca/modelo del plotter
- Plotters chinos comunes: GCC, Silhouette, Cricut, o genéricos
- Generalmente usan puerto COM virtual (USB-serial)

## Uso en el sistema
- El sistema NO controla el plotter directamente
- El admin ve la cola de corte en el panel web
- Abre el mismo PDF que usó para imprimir
- Configura el corte en el software del plotter
- Marca "Cortado" en el sistema

## Marcas de registro
- El diseño debe incluir marcas de registro (crop marks)
- El plotter las usa para alinear el corte con la impresión
- Sin marcas de registro, el corte será impreciso

## Materiales cortables
| Material | Configuración |
|----------|--------------|
| Vinilo adhesivo | Presión media, cuchilla nueva |
| Papel adhesivo | Presión baja, cuchilla estándar |
| Cartulina | Presión alta, cuchilla reforzada |

## Calibración
1. Imprimir hoja de prueba con marcas de registro
2. Colocar en el plotter
3. Ajustar offset X/Y en el software del plotter
4. Probar corte en esquina superior izquierda
5. Repetir hasta que corte exacto
