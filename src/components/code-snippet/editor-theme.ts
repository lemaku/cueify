import { EditorView } from '@codemirror/view'

const accent60 = 'var(--vt-c-accent-60)',
  accent80 = 'var(--vt-c-accent-80)',
  backgroundColor = 'var(--color-code-editor-background)',
  textColor = 'var(--color-code-editor-text)',
  lineNumberColor = 'var(--color-code-editor-linenumber)',
  accent = 'var(--vt-c-accent)',
  acccent20 = 'var(--vt-c-accent-20)',
  transparent = 'transparent'

export const editTheme = EditorView.theme(
  {
    '&': {
      backgroundColor: backgroundColor,
      color: textColor,
      borderRadius: '0.4rem'
    },
    '.cm-content': {
      caretColor: accent
    },
    '.cm-cursor, .cm-dropCursor': { borderLeftColor: accent },
    '&.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection':
      { backgroundColor: accent80 },
    '.cm-panels': { backgroundColor: transparent, color: accent60 },
    '.cm-panels.cm-panels-top': { borderBottom: '2px solid black' },
    '.cm-panels.cm-panels-bottom': { borderTop: '2px solid black' },
    '.cm-activeLine': { backgroundColor: acccent20 },
    '.cm-selectionMatch': { backgroundColor: accent60 },
    '&.cm-focused .cm-matchingBracket, &.cm-focused .cm-nonmatchingBracket': {
      backgroundColor: accent60
    },
    '.cm-gutters': {
      backgroundColor: transparent,
      color: lineNumberColor,
      border: 'none'
    },
    '.cm-activeLineGutter': {
      backgroundColor: accent60
    },
    '.cm-tooltip': {
      border: 'none',
      backgroundColor: transparent
    },
    '.cm-tooltip .cm-tooltip-arrow:before': {
      borderTopColor: transparent,
      borderBottomColor: transparent
    },
    '.cm-tooltip .cm-tooltip-arrow:after': {
      borderTopColor: transparent,
      borderBottomColor: transparent
    },
    '.cm-tooltip-autocomplete': {
      '& > ul > li[aria-selected]': {
        backgroundColor: transparent,
        color: accent60
      }
    }
  },
  { dark: false }
)

export const readOnlyTheme = EditorView.theme(
  {
    '&': {
      backgroundColor: backgroundColor,
      color: textColor,
      borderRadius: '0.4rem'
    },
    '.cm-content': {
      caretColor: accent
    },
    '.cm-cursor, .cm-dropCursor': { borderLeftColor: accent },
    '&.cm-focused > .cm-scroller > .cm-selectionLayer .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection':
      { backgroundColor: accent60 },
    '.cm-panels': { backgroundColor: transparent, color: accent60 },
    '.cm-panels.cm-panels-top': { borderBottom: '2px solid black' },
    '.cm-panels.cm-panels-bottom': { borderTop: '2px solid black' },
    '.cm-activeLine': { backgroundColor: transparent },
    '.cm-selectionMatch': { backgroundColor: transparent },
    '&.cm-focused .cm-matchingBracket, &.cm-focused .cm-nonmatchingBracket': {
      backgroundColor: transparent
    },
    '.cm-gutters': {
      backgroundColor: transparent,
      color: lineNumberColor,
      border: 'none'
    },
    '.cm-activeLineGutter': {
      backgroundColor: transparent
    },
    '.cm-tooltip': {
      border: 'none',
      backgroundColor: transparent
    },
    '.cm-tooltip .cm-tooltip-arrow:before': {
      borderTopColor: transparent,
      borderBottomColor: transparent
    },
    '.cm-tooltip .cm-tooltip-arrow:after': {
      borderTopColor: transparent,
      borderBottomColor: transparent
    },
    '.cm-tooltip-autocomplete': {
      '& > ul > li[aria-selected]': {
        backgroundColor: transparent,
        color: accent60
      }
    }
  },
  { dark: false }
)
