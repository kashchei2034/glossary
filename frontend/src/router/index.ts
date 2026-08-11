import { createRouter, createWebHistory } from 'vue-router';
import GlossaryView from '@/views/GlossaryView.vue';
import AdminView from '@/views/AdminView.vue';
import EditorView from '@/views/EditorView.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: GlossaryView,
  },
  {
    path: '/doc/:slug',
    name: 'DocumentDetail',
    component: GlossaryView,
  },
  {
    path: '/admin',
    name: 'Admin',
    component: AdminView,
  },
  {
    path: '/editor/new',
    name: 'CreateEditor',
    component: EditorView,
  },
  {
    path: '/editor/:id',
    name: 'EditEditor',
    component: EditorView,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 };
  },
});

export default router;
