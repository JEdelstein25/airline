import { env } from '$env/dynamic/public'
import type { paths } from '$lib/airline.openapi'
import createClient, { type Middleware } from 'openapi-fetch'
import { schema } from './airline.typebox'

export const apiClient = createClient<paths>({
	baseUrl: (env.API_GATEWAY_URL || `http://localhost:${env.PUBLIC_API_PORT ?? '8000'}`),
	path: env.PUBLIC_API_BASE_URL || '/api',

	// In server-side rendering context, we need the full URL to the API from the Docker container
	initFetch: (url, init) => {
		// During SSR, change the host to the API gateway service name for container communication
		if (typeof window === 'undefined') {
			// Use Docker service name instead of localhost for SSR requests
			url = url.replace(/http:\/\/localhost:8000|http:\/\/127.0.0.1:8000/, 'http://api-gateway:8000');
			console.log('SSR request URL:', url);
		} else {
			// For client-side requests from the browser, ensure we use localhost
			url = url.replace('http://api-gateway:8000', 'http://localhost:8000');
			console.log('Client request URL:', url);
		}
		return fetch(url, init);
	},
})

const detectResponseError: Middleware = {
	async onResponse({ response }) {
		if (!response.ok) {
			// TODO!(sqs): does this work?
			throw new Error(`Error ${response.status}: ${response.statusText}`)
		}
		return response
	},
}
apiClient.use(detectResponseError)

// Workaround for the sveltekit-superforms error "Multi-type unions must have a default value, or exactly one of the union types must have.".
// components['schemas']['AirlineSpec'].default = undefined
// components['schemas']['AircraftSpec'].default = undefined
export function workaroundForMultiTypeUnions(): void {
	schema['/aircraft'].POST.args.properties.body.properties.airline.default = undefined
	schema['/aircraft/{aircraftSpec}'].PATCH.args.properties.body.properties.airline.default =
		undefined
	schema['/flights'].POST.args.properties.body.properties.fleet.default = undefined
	schema['/flights/{id}'].PATCH.args.properties.body.properties.fleet.default = undefined
	schema['/flights'].POST.args.properties.body.properties.aircraft.default = undefined
	schema['/flights/{id}'].PATCH.args.properties.body.properties.aircraft.default =
		undefined
	schema['/schedules'].POST.args.properties.body.properties.airline.default = undefined
	schema['/schedules'].POST.args.properties.body.properties.originAirport.default = undefined
	schema['/schedules'].POST.args.properties.body.properties.destinationAirport.default =
		undefined
	schema['/schedules'].POST.args.properties.body.properties.fleet.default = undefined
	schema['/schedules/{id}'].PATCH.args.properties.body.properties.originAirport.default =
		undefined
	schema[
		'/schedules/{id}'
	].PATCH.args.properties.body.properties.destinationAirport.default = undefined
	schema['/schedules/{id}'].PATCH.args.properties.body.properties.fleet.default = undefined
}
