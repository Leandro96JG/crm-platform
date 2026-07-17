import { fetchPrintJobs } from './search';
import { updatePrintJobStatus } from './update-status';

export const crmCoreEndpoint = process.env.CRM_CORE_ENDPOINT;
export const crmCoreApiKey = process.env.CRM_CORE_API_KEY;

export { fetchPrintJobs, updatePrintJobStatus };
