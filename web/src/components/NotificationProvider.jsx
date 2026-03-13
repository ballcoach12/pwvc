import { Alert, Snackbar } from '@mui/material'
import { createContext, useContext, useState } from 'react'

const NotificationContext = createContext()

export const useSnackbar = () => {
  const context = useContext(NotificationContext)
  if (!context) {
    throw new Error('useSnackbar must be used within a NotificationProvider')
  }
  return context
}

export const NotificationProvider = ({ children }) => {
  const [notifications, setNotifications] = useState([])

  const enqueueSnackbar = (message, options = {}) => {
    const { variant = 'default', autoHideDuration = 5000 } = options
    const id = Date.now() + Math.random()
    
    const notification = {
      id,
      message,
      variant,
      open: true,
    }

    setNotifications(prev => [...prev, notification])

    // Auto hide
    setTimeout(() => {
      setNotifications(prev => prev.filter(n => n.id !== id))
    }, autoHideDuration)

    return id
  }

  const closeSnackbar = (id) => {
    setNotifications(prev => prev.filter(n => n.id !== id))
  }

  return (
    <NotificationContext.Provider value={{ enqueueSnackbar, closeSnackbar }}>
      {children}
      {notifications.map(notification => (
        <Snackbar
          key={notification.id}
          open={notification.open}
          onClose={() => closeSnackbar(notification.id)}
          anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
        >
          <Alert
            onClose={() => closeSnackbar(notification.id)}
            severity={notification.variant}
            variant="filled"
          >
            {notification.message}
          </Alert>
        </Snackbar>
      ))}
    </NotificationContext.Provider>
  )
}