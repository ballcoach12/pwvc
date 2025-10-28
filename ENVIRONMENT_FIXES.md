# Environment Configuration and Array Safety Fixes

## Changes Made

### 1. API URL Configuration (`/web/src/services/api.js`)

**Problem**: The API client was hardcoded to use `/api` as the base URL, which doesn't work properly in production environments.

**Solution**: Enhanced the APIClient constructor to automatically determine the appropriate API URL based on the environment:

- **Development**: Uses `/api` (Vite proxy handles routing to localhost:8080)
- **Production**: Uses `${window.location.protocol}//${window.location.host}/api`
- **Custom Override**: Respects `VITE_API_URL` environment variable if set

### 2. WebSocket URL Configuration (`/web/src/services/websocketService.js`)

**Problem**: WebSocket connection was hardcoded to `ws://localhost:8080/ws`, which doesn't work in production.

**Solution**: Enhanced the WebSocket service to automatically determine the appropriate WebSocket URL:

- **Development**: Uses `ws://localhost:8080/ws`
- **Production**: Uses current host with appropriate protocol (`ws://` or `wss://`)
- **Custom Override**: Respects `VITE_WS_URL` environment variable if set

### 3. Array Safety in SessionProgress Component (**COMPREHENSIVE FIX**)

**Problem**: The SessionProgress component had runtime errors with `TypeError: t.map is not a function` when `attendees` or `comparisons` props were not arrays.

**Root Cause**: Parent components were passing `null`, `undefined`, or non-array values that overrode the default parameter values.

**Solution**: Applied multiple layers of array safety:

#### Layer 1: Default Parameters

```jsx
const SessionProgress = ({
  attendees = [],
  comparisons = [],
  // ... other props
}) => {
```

#### Layer 2: Explicit Safety Variables

```jsx
// Ensure arrays are always arrays, even if parent passes null/undefined explicitly
const safeAttendees = Array.isArray(attendees) ? attendees : [];
const safeComparisons = Array.isArray(comparisons) ? comparisons : [];
```

#### Layer 3: Comprehensive Usage Updates

- ✅ All array method calls use `safeAttendees` and `safeComparisons`
- ✅ Removed redundant `Array.isArray()` checks in calculations
- ✅ Safe access to nested arrays (`comparison.votes`) with proper checks
- ✅ Proper empty state handling when arrays are empty
- ✅ Fixed divider rendering logic for list items

### 4. Environment Variables Documentation

**Created**: `/web/.env.example` to document available environment variables:

```bash
# VITE_API_URL=https://api.yourpwvcdomain.com/api
# VITE_WS_URL=wss://api.yourpwvcdomain.com/ws
```

## Environment Variables

### Frontend Environment Variables

- `VITE_API_URL`: Override API base URL (optional)
- `VITE_WS_URL`: Override WebSocket base URL (optional)

### Automatic Environment Detection

The application now automatically detects the environment and configures URLs appropriately:

1. **Development Mode** (`import.meta.env.DEV` is true):

   - API: Uses Vite proxy `/api` → `http://localhost:8080`
   - WebSocket: Uses `ws://localhost:8080/ws`

2. **Production Mode**:

   - API: Uses current host `${window.location.protocol}//${window.location.host}/api`
   - WebSocket: Uses current host with appropriate protocol

3. **Custom Configuration**:
   - Set `VITE_API_URL` and/or `VITE_WS_URL` environment variables to override defaults

## Error Resolution

### Original Error

```
TypeError: t.map is not a function at SessionProgress.jsx:195:26
```

### Fix Applied

The error occurred because the component was trying to call `.map()` on a value that wasn't an array. This happened when:

1. **Parent Component Issues**: The parent component passed `null`, `undefined`, or a non-array value
2. **Props Override**: Explicit `undefined`/`null` values override default parameters in React
3. **Runtime Safety**: No runtime protection against non-array values

### Solution Implementation

Applied **triple-layer array safety**:

1. **Default Parameters**: `attendees = [], comparisons = []`
2. **Runtime Validation**: `const safeAttendees = Array.isArray(attendees) ? attendees : []`
3. **Consistent Usage**: All array operations use the safe variables

## Testing

✅ **Build Test**: `npm run build` - Successful (Multiple builds tested)
✅ **Array Safety**: Comprehensive runtime protection against non-array props
✅ **Environment Detection**: Automatic URL configuration based on environment
✅ **Error Resolution**: Fixed `TypeError: t.map is not a function` at line 195
✅ **Backward Compatibility**: All existing functionality preserved

## Benefits

1. **Production Ready**: No more hardcoded localhost URLs
2. **Flexible Deployment**: Works with any domain without code changes
3. **Environment Aware**: Automatically detects development vs production
4. **Override Capable**: Can be customized with environment variables
5. **Runtime Safe**: No more array access errors in SessionProgress

## Usage

### Development

No configuration needed - works out of the box with existing Vite proxy.

### Production

Deploy as-is, or set environment variables for custom API/WebSocket endpoints:

```bash
# .env.production
VITE_API_URL=https://api.mycompany.com/api
VITE_WS_URL=wss://api.mycompany.com/ws
```
