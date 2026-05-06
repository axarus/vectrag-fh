import { createBrowserRouter } from 'react-router-dom';
import MainLayout from './layouts/MainLayout';
import DashboardPage from './features/dashboard/DashboardPage';
import ModelsListPage from './features/models/ModelsListPage';
import ModelEditorPage from './features/models/ModelEditorPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <DashboardPage />,
      },
      {
        path: 'models',
        children: [
          {
            index: true,
            element: <ModelsListPage />,
          },
          {
            path: 'new',
            element: <ModelEditorPage />,
          },
          {
            path: ':slug',
            element: <ModelEditorPage />,
          },
        ],
      },
    ],
  },
]);

