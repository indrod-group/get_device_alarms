# Informe de Software

## Introducción
El software que se ha desarrollado es un sistema de seguimiento de alarmas. Este sistema se encarga de obtener información de los dispositivos, generar solicitudes de alarmas, ejecutar estas solicitudes y finalmente guardar los datos de las alarmas. El sistema se ejecuta en ciclos, con cada ciclo iniciándose cada minuto.

## Diseño y Patrones de Programación
El diseño del software se basa en el patrón de diseño de la Cadena de Responsabilidad. Este patrón de diseño permite que una solicitud pase a través de una cadena de manejadores. Cada manejador decide si puede manejar la solicitud o si debe pasarla al siguiente manejador en la cadena.

En este caso, los manejadores son `DeviceController`, `RequestGenerator`, `RequestExecutor` y `DataSaver`. Cada uno de estos manejadores implementa la interfaz `Handler`, que define dos métodos: `Handle` y `SetNext`.

### DeviceController
`DeviceController` es el primer manejador en la cadena. Su tarea es obtener información de los dispositivos. Para hacer esto, realiza una solicitud HTTP a una API y decodifica la respuesta en una lista de dispositivos. Si ocurre un error durante este proceso, `DeviceController` utiliza la lista de dispositivos obtenida en la última solicitud exitosa.

### RequestGenerator
`RequestGenerator` es el segundo manejador en la cadena. Su tarea es generar las URLs que se utilizarán para las solicitudes de alarmas. Para hacer esto, toma la lista de dispositivos obtenida por `DeviceController` y genera una URL para cada dispositivo.

### RequestExecutor
`RequestExecutor` es el tercer manejador en la cadena. Su tarea es ejecutar las solicitudes de alarmas. Para hacer esto, toma las URLs generadas por `RequestGenerator` y realiza una solicitud HTTP a cada URL. Luego decodifica la respuesta de cada solicitud en un objeto `AlarmResponse`.

### DataSaver
`DataSaver` es el último manejador en la cadena. Su tarea es guardar los datos de las alarmas. Para hacer esto, toma los objetos `AlarmResponse` obtenidos por `RequestExecutor` y los convierte en objetos `Alarm`. Luego, realiza una solicitud HTTP para cada objeto `Alarm` para guardar los datos de la alarma.

## Director
El `Director` es responsable de construir la cadena de manejadores y procesar las solicitudes. Utiliza el patrón de diseño Singleton para asegurarse de que solo exista una instancia de `Director` en el programa. El `Director` construye la cadena de manejadores en el método `BuildChain` y procesa las solicitudes en el método `ProcessRequest`.

## Conclusión
El sistema de seguimiento de alarmas es un ejemplo de cómo se pueden utilizar los patrones de diseño de la Cadena de Responsabilidad y Singleton para crear un sistema robusto y bien estructurado. Este sistema es capaz de manejar errores y recuperarse de ellos, y puede realizar múltiples tareas en paralelo para mejorar la eficiencia. Aunque el sistema es complejo, su diseño modular hace que sea fácil de entender y mantener.
