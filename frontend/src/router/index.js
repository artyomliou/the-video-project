import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import PasswordView from '../views/PasswordView.vue'
import index from '../views/index.vue'
import notFound from '../views/notFound.vue'


const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'index',
            alias: ['/home, /index'],
            component: index
        },
        {
            path: '/login',
            name: 'login',
            component: LoginView
        },
        {
            path: '/password',
            name: 'password',
            component: PasswordView
        },
        { 
            // form Kuro https://book.vue.tw/CH4/4-2-route-settings.html
            path: '/:pathMatch(.*)*',
            name: 'notFound',
            component: notFound
        },
    ]
});

// 每一個 router 都會經過
// router.beforeEach((to) => {
//     const store = useLoginCheckStore()
//     if(to.meta.requiresAuth && !store.isLogined)
//     return '/login'
// })

// router.beforeEach( async(to) => {
//     //如果已登入(isLogined)就轉到首頁
//     const store = useLoginCheckStore()
//     if(
//         store.isLogined && //檢查登入
//         to.name !== 'login'
//     ) {
//         return { name: 'index'}
//     }
// });
// // note: meta https://router.vuejs.org/zh/guide/advanced/meta.html


export default router