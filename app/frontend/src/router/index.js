import {createRouter, createWebHistory} from 'vue-router'
import Api from '../views/Api.vue'
import {renderIcon} from "../utils/common";
import SendAlt from "../assets/icons/SendAlt.svg";

export const routes = [
    {path: '/', redirect: '/api'},
    {
        path: '/api', name: 'Api', component: Api, meta: {
            label: 'API客户端',
            icon: renderIcon(SendAlt),
        }
    },
   
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router