# Documentación del Servidor "UbicaBus"

## Índice

1. [Introducción](#introducción)
2. [Estructura del Proyecto](#estructura-del-proyecto)
3. [Requisitos y Dependencias](#requisitos-y-dependencias)
4. [Instalación y Configuración](#instalación-y-configuración)
5. [Ejecutando el Servidor](#ejecutando-el-servidor)
6. [Endpoints y Funcionalidades](#endpoints-y-funcionalidades)
   - [HTTP con Gin](#http-con-gin)
   - [Servidor WebSockets](#servidor-websockets)
   - [Cliente MQTT](#cliente-mqtt)
7. [Arquitectura y Diseño](#arquitectura-y-diseño)
8. [Formato de los Commits](#formato-de-los-commits)
9. [Consideraciones Finales](#consideraciones-finales)
10. [Referencias](#referencias)

---

## Introducción

El proyecto **UbicaBus** es un ejemplo de servidor desarrollado en Go que implementa una arquitectura limpia y modular. Se han integrado tres tecnologías clave para la comunicación:

- **HTTP**: Utilizando el framework Gin para exponer endpoints REST.
- **WebSockets**: Empleando la librería Gorilla WebSocket para crear un servidor de WebSockets.
- **MQTT**: Integrando Eclipse Paho para gestionar mensajes vía MQTT.

Este enfoque permite construir un backend escalable, de fácil mantenimiento y adaptable a múltiples canales de comunicación.

---

## Estructura del Proyecto

La organización del código se ha dividido en capas, siguiendo un patrón de separación de responsabilidades:

```
/UbicaBus
 ├── cmd
 │    └── main.go                // Punto de entrada de la aplicación
 ├── domain
 │    └── user.go                // Ejemplo de entidad (modelo de dominio)
 ├── application
 │    └── user_service.go        // Lógica de negocio (servicios de aplicación)
 ├── infrastructure
 │    ├── delivery
 │    │    ├── http_handler.go   // Endpoints HTTP con Gin
 │    │    ├── mqtt_handler.go   // Configuración y manejo del cliente MQTT
 │    │    └── websocket_handler.go  // Servidor WebSockets
 │    └── persistence
 │         └── db.go             // Conexión a la base de datos (MongoDB, patrón Singleton)
 └── go.mod                     // Gestión de dependencias del proyecto
```

Cada capa cumple un propósito específico:

- **cmd:** Archivo principal que inicializa la aplicación.
- **domain:** Definición de entidades y modelos de negocio.
- **application:** Contiene la lógica de negocio.
- **infrastructure:** Implementa detalles técnicos como la entrega (HTTP, WebSockets, MQTT) y la persistencia (MongoDB).

---

## Requisitos y Dependencias

Antes de ejecutar el servidor, asegúrate de contar con:

- **Go** (versión 1.20 o superior).
- **MongoDB** (ejecutándose localmente o accesible vía URI).
- **Broker MQTT** (ejecutándose localmente en `tcp://localhost:1883` o ajusta la configuración).

Dependencias instaladas vía `go get`:

- Gin: `github.com/gin-gonic/gin`
- MongoDB Driver: `go.mongodb.org/mongo-driver/mongo`
- Eclipse Paho MQTT: `github.com/eclipse/paho.mqtt.golang`
- Gorilla WebSocket: `github.com/gorilla/websocket`

---

## Instalación y Configuración

1. **Inicializa el módulo Go (si aún no lo has hecho):**

   ```bash
   go mod init UbicaBus
   ```

2. **Descarga las dependencias necesarias:**

   ```bash
   go get github.com/gin-gonic/gin
   go get go.mongodb.org/mongo-driver/mongo
   go get github.com/eclipse/paho.mqtt.golang
   go get github.com/gorilla/websocket
   ```

3. **Configura la conexión a MongoDB:**  
   En `infrastructure/persistence/db.go` se establece la conexión a MongoDB utilizando el patrón Singleton. Verifica que la URI `"mongodb://localhost:27017"` corresponda a tu configuración.

4. **Configura el Broker MQTT:**  
   En `infrastructure/delivery/mqtt_handler.go` se configura la conexión y suscripción al broker MQTT. Asegúrate de que el broker esté corriendo y que la URL (`tcp://localhost:1883`) sea la correcta.

---

## Ejecutando el Servidor

Para iniciar el servidor, desde la raíz del proyecto ejecuta:

```bash
go run cmd/main.go
```

El servidor iniciará:
- El **HTTP Server** en el puerto `8080`.
- El **Cliente MQTT** en una goroutine, manteniendo la conexión al broker.
- El **Servidor WebSockets** integrado en el endpoint HTTP.

---

## Endpoints y Funcionalidades

### HTTP con Gin

El servidor expone los siguientes endpoints:

- **GET /hola**  
  Responde con un mensaje JSON: `"Hola Mundo"`.

- **GET /users**  
  Devuelve una lista simulada de usuarios utilizando el servicio definido en `application/user_service.go`.

### Servidor WebSockets

- **GET /ws**  
  Establece una conexión WebSocket mediante Gorilla WebSocket.  
  - **Funcionamiento:**  
    Una vez establecida la conexión, el servidor actúa como un "echo server", devolviendo cualquier mensaje recibido.

### Cliente MQTT

El servicio MQTT se ejecuta en paralelo al servidor HTTP:
- Se conecta al broker configurado (`tcp://localhost:1883`).
- Se suscribe al tópico `"mi/topico"`.
- Los mensajes recibidos en ese tópico se registran en la consola del servidor.

---

## Arquitectura y Diseño

El proyecto sigue los principios de la **arquitectura limpia** y la **arquitectura hexagonal**:

- **Separación de Responsabilidades:**  
  Cada capa (dominio, aplicación, infraestructura) está aislada, permitiendo modificar o extender funcionalidades sin afectar el núcleo del negocio.

- **Patrón Singleton para la Base de Datos:**  
  Se utiliza en `infrastructure/persistence/db.go` para garantizar una única instancia de conexión a MongoDB.

- **Integración de Múltiples Protocolos:**  
  El servidor combina HTTP, WebSockets y MQTT, facilitando la comunicación a través de diferentes canales.

---

## Formato de los Commits

Para mantener un historial de cambios claro y consistente, se recomienda seguir el formato **Conventional Commits**. Este formato facilita la generación automática de logs, changelogs y versiones semánticas. Algunos ejemplos de formatos de commit son:

- **feat:** Para nuevas funcionalidades.  
  Ejemplo:  
  ```
  feat: agregar endpoint para obtener usuarios
  ```
- **fix:** Para correcciones de errores.  
  Ejemplo:  
  ```
  fix: corregir error en conexión a MongoDB
  ```
- **docs:** Para cambios en la documentación.  
  Ejemplo:  
  ```
  docs: actualizar documentación de endpoints
  ```
- **chore:** Para tareas de mantenimiento, configuraciones o actualizaciones de dependencias.  
  Ejemplo:  
  ```
  chore: actualizar versión de Gin
  ```

**Recomendación en Visual Studio Code:**  
Para facilitar la adopción de este formato, se recomienda instalar la extensión **Conventional Commits** en Visual Studio Code. Esta extensión te ayudará a:
- Generar mensajes de commit siguiendo el estándar.
- Validar que los commits cumplan con el formato propuesto.
- Aumentar la consistencia y calidad de los commits en el equipo de desarrollo.

Puedes instalar la extensión buscando "Conventional Commits" en el Marketplace de Visual Studio Code o visitando el siguiente enlace: [Conventional Commits VSCode Extension](https://marketplace.visualstudio.com/items?itemName=vivaxy.vscode-conventional-commits).

---

## Consideraciones Finales

- **Escalabilidad y Mantenibilidad:**  
  La modularidad del proyecto facilita la incorporación de nuevos endpoints, servicios o la integración de otros protocolos.
  
- **Pruebas y Desarrollo:**  
  La separación de responsabilidades permite realizar pruebas unitarias en cada capa y facilita la colaboración en equipo.

- **Configuración de Producción:**  
  Antes de desplegar en producción, revisa las configuraciones (URI de MongoDB, dirección del broker MQTT, manejo de CORS en WebSockets, etc.) para asegurar un funcionamiento seguro y eficiente.

---

## Referencias

- [Refactoring Guru – Arquitectura limpia y patrones de diseño](https://refactoring.guru/es) citeturn0search0
- [Patterns.dev – Recursos y ejemplos de patrones de diseño](https://www.patterns.dev/) citeturn0search0
- [Gin Web Framework](https://github.com/gin-gonic/gin) citeturn0search0
- [Gorilla WebSocket](https://github.com/gorilla/websocket) citeturn0search0
- [Eclipse Paho MQTT Client](https://github.com/eclipse/paho.mqtt.golang) citeturn0search0
- [Conventional Commits Extension for VSCode](https://marketplace.visualstudio.com/items?itemName=vivaxy.vscode-conventional-commits) citeturn0search0