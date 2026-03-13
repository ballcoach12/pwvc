import { Keyboard } from '@mui/icons-material'
import {
    Box,
    Button,
    Chip,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle,
    Divider,
    Grid,
    List,
    ListItem,
    ListItemText,
    Typography
} from '@mui/material'
import React from 'react'
import { useShortcutHelp } from '../hooks/useKeyboardShortcuts'

/**
 * KeyboardShortcutsDialog displays all available keyboard shortcuts
 * in an organized, easy-to-read format with visual key indicators.
 */
const KeyboardShortcutsDialog = ({ open, onClose }) => {
  const { getShortcutHelp } = useShortcutHelp()
  const shortcutHelp = getShortcutHelp()

  const formatKey = (key) => {
    const keyMap = {
      ' ': 'Space',
      'ArrowLeft': '←',
      'ArrowRight': '→',
      'ArrowUp': '↑',
      'ArrowDown': '↓',
      'Escape': 'Esc'
    }
    return keyMap[key] || key
  }

  const renderKeyChip = (key) => (
    <Chip
      key={key}
      label={formatKey(key)}
      size="small"
      variant="outlined"
      sx={{
        fontFamily: 'monospace',
        fontWeight: 'bold',
        mr: 0.5
      }}
    />
  )

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="md"
      fullWidth
      PaperProps={{
        sx: { minHeight: '60vh' }
      }}
    >
      <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Keyboard />
        Keyboard Shortcuts
      </DialogTitle>
      
      <DialogContent>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
          Use these keyboard shortcuts to navigate and vote more efficiently during pairwise comparison sessions.
        </Typography>

        <Grid container spacing={3}>
          {shortcutHelp.map((category, categoryIndex) => (
            <Grid item xs={12} md={4} key={category.category}>
              <Typography variant="h6" gutterBottom color="primary">
                {category.category}
              </Typography>
              
              <List dense sx={{ bgcolor: 'grey.50', borderRadius: 1, p: 1 }}>
                {category.shortcuts.map((shortcut, index) => (
                  <React.Fragment key={index}>
                    <ListItem sx={{ px: 1, py: 0.5 }}>
                      <ListItemText
                        primary={
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                            {shortcut.keys.map(renderKeyChip)}
                          </Box>
                        }
                        secondary={
                          <Typography variant="body2" color="text.secondary">
                            {shortcut.description}
                          </Typography>
                        }
                      />
                    </ListItem>
                    {index < category.shortcuts.length - 1 && (
                      <Divider variant="inset" sx={{ ml: 1, mr: 1 }} />
                    )}
                  </React.Fragment>
                ))}
              </List>
            </Grid>
          ))}
        </Grid>

        <Box sx={{ mt: 4, p: 2, bgcolor: 'info.light', borderRadius: 1 }}>
          <Typography variant="body2" sx={{ fontWeight: 'medium', mb: 1 }}>
            💡 Pro Tips:
          </Typography>
          <Typography variant="body2" component="div">
            <ul style={{ margin: 0, paddingLeft: '1.2em' }}>
              <li>Keyboard shortcuts work in both grid and detail view modes</li>
              <li>Shortcuts are disabled when typing in input fields</li>
              <li>Use fullscreen mode (F key) for distraction-free voting</li>
              <li>Auto-advance can be enabled in settings to move to the next comparison after voting</li>
              <li>Press Escape to exit fullscreen or close dialogs quickly</li>
            </ul>
          </Typography>
        </Box>
      </DialogContent>
      
      <DialogActions>
        <Button onClick={onClose} variant="contained">
          Got it!
        </Button>
      </DialogActions>
    </Dialog>
  )
}

export default KeyboardShortcutsDialog