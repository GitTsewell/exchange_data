import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

const login = r => require.ensure([], () => r(require('@/page/login')), 'login');
const manage = r => require.ensure([], () => r(require('@/page/manage')), 'manage');
const depth = r => require.ensure([], () => r(require('@/page/depth')), 'depth');
const system = r => require.ensure([], () => r(require('@/page/system')), 'system');

const routes = [
	{
		path: '/',
		component: login
	},
	{
		path: '/manage',
		component: manage,
		name: '',
		children: [{
            path: '/depth',
            component: depth,
            meta: {
                requireAuth: true,  // 该路由项需要权限校验
            }
        },{
            path: '/system',
            component: system,
            meta: {
                requireAuth: true,  // 该路由项需要权限校验
            }
        }]
	}
]

export default new Router({
	routes,
	strict: process.env.NODE_ENV !== 'production',
})
