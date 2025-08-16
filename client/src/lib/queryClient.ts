import { QueryClient } from '@tanstack/react-query';

// Create a client
export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            // Time in milliseconds that unused/inactive cache data remains in memory
            gcTime: 1000 * 60 * 5, // 5 minutes
            // Time in milliseconds after data is considered stale
            staleTime: 1000 * 60 * 2, // 2 minutes
            // Retry failed requests
            retry: 2,
            // Refetch on window focus
            refetchOnWindowFocus: false,
            // Refetch on reconnect
            refetchOnReconnect: true,
        },
        mutations: {
            // Retry failed mutations
            retry: 1,
        },
    },
});