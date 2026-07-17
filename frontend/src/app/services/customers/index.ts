import { fetchCustomers } from './search';
import { fetchCustomer } from './fetch-one';
import { createCustomer } from './create';
import { updateCustomer } from './update';
import { deleteCustomer } from './delete';

export const crmCoreEndpoint = process.env.CRM_CORE_ENDPOINT;
export const crmCoreApiKey = process.env.CRM_CORE_API_KEY;

export {
  fetchCustomers,
  fetchCustomer,
  createCustomer,
  updateCustomer,
  deleteCustomer,
};
