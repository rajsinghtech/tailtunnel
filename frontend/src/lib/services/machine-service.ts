import APIService from './api-service';
import type { MachineListResponse } from '$lib/types/machine';

export default class MachineService extends APIService {
	list = async (): Promise<MachineListResponse> => {
		const res = await this.api.get('/machines');
		return res.data as MachineListResponse;
	};
}
