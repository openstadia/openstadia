export async function getApps() {
    const response = await fetch('/api/apps', {
        method: 'GET',
        headers: {
            "Content-Type": "application/json",
        }
    })

    return response.json()
}