import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { useRouter } from 'next/navigation';
import { signOut } from 'next-auth/react';
import { User, UserRole } from '../../../types/user';
import { ProfileForm } from '../profile-form';

jest.mock('next/navigation', () => ({
  useRouter: jest.fn(),
}));

jest.mock('next-auth/react', () => ({
  signOut: jest.fn(),
}));

jest.mock('../../../services/user', () => ({
  updateUser: jest.fn(),
  changePassword: jest.fn(),
  isEmailTaken: jest.fn(),
}));

const showSnackbar = jest.fn();
jest.mock('../../../context/SnackbarProvider', () => ({
  useSnackbar: () => ({ showSnackbar }),
}));

const { updateUser, changePassword, isEmailTaken } = jest.requireMock(
  '../../../services/user'
);

const mockUser: User = {
  user_id: 'user-1',
  username: 'joao.silva',
  first_name: 'João',
  last_name: 'Silva',
  email: 'joao@example.com',
  role: UserRole.OPERATOR,
  region: 1,
  created_at: '',
  updated_at: '',
  created_by: '',
  updated_by: '',
  active: true,
};

beforeEach(() => {
  jest.clearAllMocks();
  (useRouter as jest.Mock).mockReturnValue({ refresh: jest.fn() });
  isEmailTaken.mockResolvedValue(false);
  updateUser.mockResolvedValue({ success: true, message: 'ok' });
  changePassword.mockResolvedValue({ success: true, message: 'ok' });
});

function enterEditMode() {
  fireEvent.click(screen.getByRole('button', { name: 'Editar' }));
}

describe('ProfileForm', () => {
  describe('rendering', () => {
    it('should render read-only user data by default, including username and role', () => {
      render(<ProfileForm user={mockUser} />);

      expect(screen.getByText('joao.silva')).toBeInTheDocument();
      expect(screen.getByText('Operador')).toBeInTheDocument();
      expect(screen.getByText('João')).toBeInTheDocument();
      expect(screen.getByText('Silva')).toBeInTheDocument();
      expect(screen.getByText('joao@example.com')).toBeInTheDocument();
      expect(screen.queryByLabelText('Nombre')).not.toBeInTheDocument();
    });

    it('should switch to edit mode when the Editar button is clicked', () => {
      render(<ProfileForm user={mockUser} />);

      enterEditMode();

      expect(screen.getByLabelText('Nombre')).toHaveValue('João');
      expect(screen.getByLabelText('Apellido')).toHaveValue('Silva');
      expect(screen.getByLabelText('Correo electrónico')).toHaveValue('joao@example.com');
    });

    it('should return to read-only mode when Cancelar is clicked', () => {
      render(<ProfileForm user={mockUser} />);

      enterEditMode();
      fireEvent.click(screen.getByRole('button', { name: 'Cancelar' }));

      expect(screen.queryByLabelText('Nombre')).not.toBeInTheDocument();
      expect(
        screen.getByRole('button', { name: 'Editar' })
      ).toBeInTheDocument();
    });
  });

  describe('updating profile info', () => {
    it('should submit updated first name, last name and email', async () => {
      render(<ProfileForm user={mockUser} />);

      enterEditMode();
      fireEvent.change(screen.getByLabelText('Nombre'), {
        target: { value: 'Joaquim' },
      });
      fireEvent.click(screen.getByRole('button', { name: 'Guardar datos' }));

      await waitFor(() => {
        expect(updateUser).toHaveBeenCalledWith('user-1', {
          first_name: 'Joaquim',
          last_name: 'Silva',
          email: 'joao@example.com',
          updated_by: 'joao.silva',
        });
      });
      expect(showSnackbar).toHaveBeenCalledWith(
        'Datos actualizados con éxito',
        'success'
      );
    });

    it('should block submission when the new email is already taken', async () => {
      isEmailTaken.mockResolvedValue(true);
      render(<ProfileForm user={mockUser} />);

      enterEditMode();
      fireEvent.change(screen.getByLabelText('Correo electrónico'), {
        target: { value: 'other@example.com' },
      });
      fireEvent.click(screen.getByRole('button', { name: 'Guardar datos' }));

      await waitFor(() => {
        expect(
          screen.getByText('Este email já está em uso por outro usuário')
        ).toBeInTheDocument();
      });
      expect(updateUser).not.toHaveBeenCalled();
    });

    it('should sign out when the update request is unauthorized', async () => {
      updateUser.mockResolvedValue({
        success: false,
        message: 'usuário não autorizado',
        unauthorized: true,
      });
      render(<ProfileForm user={mockUser} />);

      enterEditMode();
      fireEvent.click(screen.getByRole('button', { name: 'Guardar datos' }));

      await waitFor(() => {
        expect(signOut).toHaveBeenCalledWith({ callbackUrl: '/login' });
      });
    });
  });

  describe('changing password', () => {
    function fillPasswordForm({
      current = 'oldPass123!',
      next = 'NewPass123!',
      confirm = 'NewPass123!',
    }: { current?: string; next?: string; confirm?: string } = {}) {
      fireEvent.change(screen.getByLabelText('Contraseña actual'), {
        target: { value: current },
      });
      fireEvent.change(screen.getByLabelText('Nueva contraseña'), {
        target: { value: next },
      });
      fireEvent.change(screen.getByLabelText('Confirme la nueva contraseña'), {
        target: { value: confirm },
      });
      fireEvent.click(screen.getByRole('button', { name: 'Cambiar contraseña' }));
    }

    it('should submit the current and new password', async () => {
      render(<ProfileForm user={mockUser} />);

      fillPasswordForm();

      await waitFor(() => {
        expect(changePassword).toHaveBeenCalledWith('user-1', {
          old_password: 'oldPass123!',
          new_password: 'NewPass123!',
        });
      });
      expect(showSnackbar).toHaveBeenCalledWith(
        'Contraseña actualizada con éxito',
        'success'
      );
    });

    it('should reject a new password shorter than 8 characters', async () => {
      render(<ProfileForm user={mockUser} />);

      fillPasswordForm({ next: 'Ab1!', confirm: 'Ab1!' });

      await waitFor(() => {
        expect(
          screen.getByText('La contraseña debe tener al menos 8 caracteres')
        ).toBeInTheDocument();
      });
      expect(changePassword).not.toHaveBeenCalled();
    });

    it('should reject a new password without a number', async () => {
      render(<ProfileForm user={mockUser} />);

      fillPasswordForm({ next: 'Abcdefgh!', confirm: 'Abcdefgh!' });

      await waitFor(() => {
        expect(
          screen.getByText('La contraseña debe contener al menos 1 número')
        ).toBeInTheDocument();
      });
      expect(changePassword).not.toHaveBeenCalled();
    });

    it('should reject a new password without a special character', async () => {
      render(<ProfileForm user={mockUser} />);

      fillPasswordForm({ next: 'Abcdefg1', confirm: 'Abcdefg1' });

      await waitFor(() => {
        expect(
          screen.getByText('La contraseña debe contener al menos 1 carácter especial')
        ).toBeInTheDocument();
      });
      expect(changePassword).not.toHaveBeenCalled();
    });

    it('should reject when confirmation does not match', async () => {
      render(<ProfileForm user={mockUser} />);

      fillPasswordForm({ confirm: 'Different123!' });

      await waitFor(() => {
        expect(
          screen.getByText('La nueva contraseña y la confirmación deben coincidir')
        ).toBeInTheDocument();
      });
      expect(changePassword).not.toHaveBeenCalled();
    });

    it('should show the API error when the current password is wrong', async () => {
      changePassword.mockResolvedValue({
        success: false,
        message: 'old_password is incorrect',
      });
      render(<ProfileForm user={mockUser} />);

      fillPasswordForm();

      await waitFor(() => {
        expect(
          screen.getByText('old_password is incorrect')
        ).toBeInTheDocument();
      });
    });
  });
});
