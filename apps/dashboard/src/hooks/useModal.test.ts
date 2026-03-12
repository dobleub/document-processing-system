import { renderHook, act } from '@testing-library/react';
import { useModal } from './useModal';

describe('useModal', () => {
  it('should initialize with default closed state', () => {
    const { result } = renderHook(() => useModal());

    expect(result.current.isOpen).toBe(false);
  });

  it('should initialize with custom initial state', () => {
    const { result } = renderHook(() => useModal(true));

    expect(result.current.isOpen).toBe(true);
  });

  it('should open modal', () => {
    const { result } = renderHook(() => useModal(false));

    act(() => {
      result.current.openModal();
    });

    expect(result.current.isOpen).toBe(true);
  });

  it('should close modal', () => {
    const { result } = renderHook(() => useModal(true));

    act(() => {
      result.current.closeModal();
    });

    expect(result.current.isOpen).toBe(false);
  });

  it('should toggle modal state', () => {
    const { result } = renderHook(() => useModal(false));

    act(() => {
      result.current.toggleModal();
    });
    expect(result.current.isOpen).toBe(true);

    act(() => {
      result.current.toggleModal();
    });
    expect(result.current.isOpen).toBe(false);
  });

  it('should expose all control methods', () => {
    const { result } = renderHook(() => useModal());

    expect(typeof result.current.openModal).toBe('function');
    expect(typeof result.current.closeModal).toBe('function');
    expect(typeof result.current.toggleModal).toBe('function');
    expect(typeof result.current.isOpen).toBe('boolean');
  });
});
