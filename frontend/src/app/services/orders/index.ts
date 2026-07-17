import { fetchOrders } from './search';
import { updateOrderStatus } from './update-status';
import { createOrder } from './create';

export const crmCoreEndpoint = process.env.CRM_CORE_ENDPOINT;
export const crmCoreApiKey = process.env.CRM_CORE_API_KEY;

export { fetchOrders, updateOrderStatus, createOrder };
