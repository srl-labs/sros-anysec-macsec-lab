import { error } from '@sveltejs/kit';

export async function load({ url, fetch }) {
  try {
    const urlHost = url.hostname
    const fetchUrl = "/api/get_all_state"
    // commented for testing
    //const resp = await fetch(fetchUrl)
    //return { urlHost: urlHost, fetchUrl: fetchUrl, state: await resp.json() }
    return { 
      urlHost: urlHost, 
      fetchUrl: fetchUrl, 
      state: {
        icmp: {
          Vll:  "enabled",
          Vpls: "enabled",
          Vprn: "enabled",
        },
        link: {
          top:    "enabled",
          bottom: "enabled",
        },
        anysec: {
          vll:  "enabled",
          vpls: "enabled",
          vprn: "enabled",
        },
      }
    }
  } catch (e) {
    throw error(404, "Backend Disconnected")
  }
}