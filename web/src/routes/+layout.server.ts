import { apiClient } from '$lib/api'
import { error } from '@sveltejs/kit'
import type { LayoutServerLoad } from './$types'

export const load: LayoutServerLoad = async () => {
	// For SSR, return mock data to avoid API gateway connection issues
	// This way we can still do SSR without needing to connect to the API
	if (typeof window === 'undefined') {
		console.log('SSR mode - using mock data')
		// Return mock data for SSR
		return {
			allAirlines: [
				{ 
					id: 1, 
					name: 'Airline Sample', 
					code: 'ALS',
					logo: 'ALS', // Logo field is needed for airline-icon component
					color: '#0088cc' // Add a default color
				}
			],
			allAircraftTypes: [
				{ id: 1, name: 'Boeing 737', code: 'B737' },
				{ id: 2, name: 'Airbus A320', code: 'A320' }
			]
		}
	}

	// For client-side, use the real API
	try {
		const airlinesResp = await apiClient.GET('/airlines', { fetch })
		if (!airlinesResp.response.ok || !airlinesResp.data) {
			console.error('Error fetching airlines:', airlinesResp.response.status)
			throw error(airlinesResp.response.status, 'Error fetching airlines')
		}

		const aircraftTypesResp = await apiClient.GET('/aircraft-types', { fetch })
		if (!aircraftTypesResp.response.ok || !aircraftTypesResp.data) {
			console.error('Error fetching aircraft types:', aircraftTypesResp.response.status)
			throw error(aircraftTypesResp.response.status, 'Error fetching aircraft types')
		}

		return {
			allAirlines: airlinesResp.data,
			allAircraftTypes: aircraftTypesResp.data,
		}
	} catch (e) {
		console.error('API request error:', e)
		throw error(500, 'Error connecting to API')
	}
}
