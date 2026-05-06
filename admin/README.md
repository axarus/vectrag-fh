# VectraG Admin Panel

Panel de administración para VectraG (modo desarrollo).

## Estructura

```
src/
├── features/           # Módulos funcionales
│   ├── dashboard/    # Panel principal
│   └── models/       # Gestión de modelos
├── layouts/          # Layouts de la aplicación
│   └── MainLayout.tsx # Layout con sidebar
├── shared/           # Recursos compartidos
│   └── services/     # Servicios API
├── types/            # Tipos TypeScript
├── routes.tsx        # Configuración de rutas
└── App.tsx           # Componente raíz
```

## Desarrollo

```bash
npm install
npm run dev
```

## Variables de Entorno

Crea un archivo `.env` basado en `.env.example`:

```
VITE_API_URL=http://localhost:8080/api
```

## Características

- ✅ Dashboard con estadísticas
- ✅ Lista de modelos
- ✅ Editor de modelos (crear/editar)
- ✅ Gestión de campos
- ✅ UI responsive con DaisyUI
- ✅ Listo para conectar al backend

## Próximos Pasos

- Conectar al servidor HTTP del backend
- Implementar autenticación
- Agregar validaciones en el frontend
- Mejorar UX/UI

