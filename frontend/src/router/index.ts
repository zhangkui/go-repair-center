import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import AppShell from '../layouts/AppShell.vue'
import ApprovalCenterView from '../views/ApprovalCenterView.vue'
import DashboardView from '../views/DashboardView.vue'
import DataToolsView from '../views/DataToolsView.vue'
import DispatchCenterView from '../views/DispatchCenterView.vue'
import GenericCrudView from '../views/GenericCrudView.vue'
import InventoryControlView from '../views/InventoryControlView.vue'
import LoginView from '../views/LoginView.vue'
import OperationsView from '../views/OperationsView.vue'
import ProfileView from '../views/ProfileView.vue'
import RegisterView from '../views/RegisterView.vue'
import SessionsView from '../views/SessionsView.vue'
import { useAuthStore } from '../stores/auth'

function col(key: string, label: string) {
  return { key, label }
}

function field(key: string, label: string, type: 'text' | 'number' = 'text') {
  return { key, label, type }
}

function crudMeta(
  title: string,
  subtitle: string,
  endpoint: string,
  columns: Array<{ key: string; label: string }>,
  fields: Array<{ key: string; label: string; type?: 'text' | 'number' }>
): Record<string, unknown> {
  return {
    menu: true,
    title,
    subtitle,
    endpoint,
    columns,
    fields
  }
}

const routes: RouteRecordRaw[] = [
  { path: '/login', component: LoginView },
  { path: '/register', component: RegisterView },
  {
    path: '/',
    component: AppShell,
    children: [
      { path: '', component: DashboardView, meta: { menu: true, title: '仪表盘' } },
      { path: 'operations', component: OperationsView, meta: { menu: true, title: '运营中心' } },
      { path: 'dispatch-center', component: DispatchCenterView, meta: { menu: true, title: '派工中心' } },
      { path: 'approval-center', component: ApprovalCenterView, meta: { menu: true, title: '审批中心' } },
      { path: 'inventory-control', component: InventoryControlView, meta: { menu: true, title: '库存管控' } },
      { path: 'sessions', component: SessionsView, meta: { menu: true, title: '会话管理' } },
      { path: 'data-tools', component: DataToolsView, meta: { menu: true, title: '数据工具' } },
      { path: 'profile', component: ProfileView, meta: { title: '个人资料' } },
      {
        path: 'users',
        component: GenericCrudView,
        meta: crudMeta(
          '用户管理',
          '维护账号信息、启停状态和密码重置。',
          '/users',
          [col('id', '编号'), col('username', '用户名'), col('display_name', '显示名称'), col('status', '状态')],
          [field('username', '用户名'), field('display_name', '显示名称'), field('password', '密码'), field('phone', '手机号'), field('email', '邮箱')]
        )
      },
      {
        path: 'roles',
        component: GenericCrudView,
        meta: crudMeta(
          '角色管理',
          '维护角色编码、名称和说明。',
          '/roles',
          [col('id', '编号'), col('code', '角色编码'), col('name', '角色名称'), col('description', '说明')],
          [field('code', '角色编码'), field('name', '角色名称'), field('description', '说明')]
        )
      },
      {
        path: 'permissions',
        component: GenericCrudView,
        meta: crudMeta(
          '权限管理',
          '查看后端强校验的权限定义。',
          '/permissions',
          [col('id', '编号'), col('code', '权限编码'), col('name', '权限名称'), col('module', '模块'), col('action', '动作')],
          [field('code', '权限编码'), field('name', '权限名称'), field('module', '模块'), field('action', '动作'), field('description', '说明')]
        )
      },
      {
        path: 'customers',
        component: GenericCrudView,
        meta: crudMeta(
          '客户管理',
          '创建并维护客户基础资料。',
          '/customers',
          [col('id', '编号'), col('name', '客户名称'), col('phone', '联系电话'), col('level', '客户等级')],
          [field('name', '客户名称'), field('phone', '联系电话'), field('address', '联系地址'), field('level', '客户等级'), field('remark', '备注')]
        )
      },
      {
        path: 'devices',
        component: GenericCrudView,
        meta: crudMeta(
          '设备档案',
          '跟踪设备身份信息与当前状态。',
          '/devices',
          [col('id', '编号'), col('brand', '品牌'), col('model', '型号'), col('serial_number', '序列号'), col('status', '状态')],
          [field('customer_id', '客户编号', 'number'), field('brand', '品牌'), field('model', '型号'), field('serial_number', '序列号'), field('category', '类别'), field('status', '状态')]
        )
      },
      {
        path: 'repair-orders',
        component: GenericCrudView,
        meta: crudMeta(
          '报修预约',
          '处理受理、派工和状态流转。',
          '/repair-orders',
          [col('id', '编号'), col('order_number', '维修单号'), col('customer_id', '客户编号'), col('device_id', '设备编号'), col('status', '状态')],
          [field('customer_id', '客户编号', 'number'), field('device_id', '设备编号', 'number'), field('fault_description', '故障描述'), field('urgency', '紧急程度'), field('service_method', '服务方式'), field('status', '状态')]
        )
      },
      {
        path: 'quotations',
        component: GenericCrudView,
        meta: crudMeta(
          '检测报价',
          '创建报价记录并维护审批状态。',
          '/quotations',
          [col('id', '编号'), col('repair_order_id', '维修单编号'), col('version', '版本'), col('total_amount', '总金额'), col('status', '状态')],
          [field('repair_order_id', '维修单编号', 'number'), field('labor_fee', '人工费', 'number'), field('parts_fee', '配件费', 'number'), field('inspection_fee', '检测费', 'number'), field('other_fee', '其他费用', 'number'), field('status', '状态')]
        )
      },
      {
        path: 'repair-executions',
        component: GenericCrudView,
        meta: crudMeta(
          '维修执行',
          '跟踪工序执行与维修状态。',
          '/repair-executions',
          [col('id', '编号'), col('repair_order_id', '维修单编号'), col('technician_id', '技师编号'), col('status', '状态')],
          [field('repair_order_id', '维修单编号', 'number'), field('technician_id', '技师编号', 'number'), field('diagnosis', '检测结论'), field('status', '状态')]
        )
      },
      {
        path: 'parts',
        component: GenericCrudView,
        meta: crudMeta(
          '配件库存',
          '维护配件资料、库存数量和库存阈值。',
          '/parts',
          [col('id', '编号'), col('code', '配件编码'), col('name', '配件名称'), col('stock_quantity', '库存数量'), col('status', '状态')],
          [field('code', '配件编码'), field('name', '配件名称'), field('specification', '规格'), field('unit', '单位'), field('unit_price', '单价', 'number'), field('stock_quantity', '库存数量', 'number'), field('safety_stock', '安全库存', 'number'), field('status', '状态')]
        )
      },
      {
        path: 'warranties',
        component: GenericCrudView,
        meta: crudMeta(
          '收费质保',
          '查看质保范围、时长和状态。',
          '/warranties',
          [col('id', '编号'), col('repair_order_id', '维修单编号'), col('warranty_type', '质保类型'), col('duration_days', '时长（天）'), col('status', '状态')],
          [field('repair_order_id', '维修单编号', 'number'), field('warranty_type', '质保类型'), field('duration_days', '时长（天）', 'number'), field('status', '状态')]
        )
      },
      {
        path: 'feedbacks',
        component: GenericCrudView,
        meta: crudMeta(
          '服务回访',
          '维护回访记录、投诉和处理结果。',
          '/feedbacks',
          [col('id', '编号'), col('repair_order_id', '维修单编号'), col('method', '回访方式'), col('status', '状态')],
          [field('repair_order_id', '维修单编号', 'number'), field('method', '回访方式'), field('complaint_type', '投诉类型'), field('complaint', '投诉内容'), field('resolution', '处理结果'), field('status', '状态')]
        )
      },
      {
        path: 'audit-logs',
        component: GenericCrudView,
        meta: crudMeta(
          '审计日志',
          '查看用户和系统操作留痕。',
          '/audit-logs',
          [col('id', '编号'), col('user_id', '用户编号'), col('action', '操作'), col('resource_type', '资源类型'), col('created_at', '创建时间')],
          [field('action', '操作'), field('resource_type', '资源类型'), field('resource_id', '资源编号')]
        )
      },
      {
        path: 'system-configs',
        component: GenericCrudView,
        meta: crudMeta(
          '系统配置',
          '维护阈值、默认值和业务配置。',
          '/system-configs',
          [col('id', '编号'), col('config_key', '配置键'), col('config_value', '配置值'), col('value_type', '值类型')],
          [field('config_key', '配置键'), field('config_value', '配置值'), field('value_type', '值类型'), field('description', '说明')]
        )
      },
      {
        path: 'fault-codes',
        component: GenericCrudView,
        meta: crudMeta(
          '故障码字典',
          '维护标准化故障码目录。',
          '/fault-codes',
          [col('id', '编号'), col('code', '故障码'), col('name', '名称'), col('category', '类别'), col('status', '状态')],
          [field('code', '故障码'), field('name', '名称'), field('category', '类别'), field('description', '说明'), field('status', '状态')]
        )
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  if (to.path === '/login' || to.path === '/register') {
    return true
  }
  if (!auth.isAuthenticated) {
    return '/login'
  }
  if (!auth.user) {
    try {
      await auth.fetchMe()
    } catch {
      return '/login'
    }
  }
  return true
})

export default router
