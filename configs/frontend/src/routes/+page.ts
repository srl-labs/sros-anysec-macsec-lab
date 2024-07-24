import { error } from '@sveltejs/kit';

export async function load({ url, fetch }) {
  try {
    const urlHost = url.hostname
    const baseUrl = `http://${urlHost}:8080`
    const fetchUrl = baseUrl + "/get_all_state"
    const resp = await fetch(fetchUrl)
    return { urlHost: urlHost, baseUrl: baseUrl, fetchUrl: fetchUrl, state: await resp.json() }
  } catch(e) {
    throw error(404, "Backend Disconnected")
  }
}