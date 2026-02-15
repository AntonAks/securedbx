import { createRouter, createWebHashHistory } from 'vue-router';

const routes = [
    {
        path: '/',
        name: 'home',
        component: () => import('../views/HomeView.vue'),
    },
    {
        path: '/download',
        name: 'download',
        component: () => import('../views/DownloadView.vue'),
    },
    {
        path: '/share',
        name: 'share',
        component: () => import('../views/ShareView.vue'),
    },
    {
        path: '/faq',
        name: 'faq',
        component: () => import('../views/FaqView.vue'),
    },
    {
        path: '/about',
        name: 'about',
        component: () => import('../views/AboutView.vue'),
    },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

export default router;
