/** @type {import('./$types').PageLoad} */
export async function load({ fetch }) {
  try {
    // Client-side fetch to API
    const airlinesRes = await fetch('/api/airlines');
    const airlines = await airlinesRes.json();

    const aircraftTypesRes = await fetch('/api/aircraft-types');
    const aircraftTypes = await aircraftTypesRes.json();

    return {
      allAirlines: airlines,
      allAircraftTypes: aircraftTypes
    };
  } catch (error) {
    console.error('Failed to load data:', error);
    return {
      allAirlines: [],
      allAircraftTypes: []
    };
  }
}