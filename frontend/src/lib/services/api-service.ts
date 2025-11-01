import axios from 'axios';

export default abstract class APIService {
	protected api = axios.create({
		baseURL: '/api'
	});
}
