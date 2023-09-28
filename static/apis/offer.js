export async function sendOffer(offer) {
    const response = await fetch('/api/offer', {
        method: 'POST',
        body: JSON.stringify(offer),
        headers: {
            "Content-Type": "application/json",
        }
    })

    return response.json()
}