import { fetchPlanchas } from './search';
import { fetchMaterials } from './materials';
import { calculatePrice } from './calculate-price';

export const crmCoreEndpoint = process.env.CRM_CORE_ENDPOINT;
export const crmCoreApiKey = process.env.CRM_CORE_API_KEY;

export { fetchPlanchas, fetchMaterials, calculatePrice };
