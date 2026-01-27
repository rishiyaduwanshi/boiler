// Cloudflare Pages catch-all function
// Redirects all other paths to GitHub repo

export async function onRequest(context) {
  const url = new URL(context.request.url);
  
  // Skip /install route (handled by install.js)
  if (url.pathname === '/install') {
    return context.next();
  }
  
  // Root redirect
  if (url.pathname === '/' || url.pathname === '') {
    return Response.redirect('https://github.com/rishiyaduwanshi/boiler', 302);
  }
  
  // 404 for other paths
  return new Response('Not Found. Visit https://github.com/rishiyaduwanshi/boiler', {
    status: 404,
    headers: { 'Content-Type': 'text/plain' }
  });
}
