class Users::RegistrationsController < Devise::RegistrationsController
  respond_to :json
  before_action :authenticate_user!, only: [:update, :destroy]

  def create
    build_resource(sign_up_params)
    authorize resource  # 👈 Pundit authorization

    if resource.save
      if resource.active_for_authentication?
        sign_up(resource_name, resource)
        render json: {
          message: 'User registered successfully',
          user: user_response(resource)
        }, status: :created
      else
        expire_data_after_sign_in!
        render json: {
          message: 'Signed up but account is inactive',
          user: user_response(resource)
        }, status: :ok
      end
    else
      render json: {
        message: 'User registration failed',
        errors: resource.errors.full_messages
      }, status: :unprocessable_entity
    end
  end

  def update
    self.resource = User.find(params[:id])
    authorize resource  # 👈 Pundit authorization

    if resource.update(account_update_params)
      render json: {
        message: 'User updated successfully',
        user: user_response(resource)
      }, status: :ok
    else
      render json: {
        message: 'User update failed',
        errors: resource.errors.full_messages
      }, status: :unprocessable_entity
    end
  end

  def destroy
    self.resource = User.find(params[:id])
    authorize resource  # 👈 Pundit authorization

    resource.destroy
    render json: { message: 'Account deleted successfully' }, status: :ok
  end

  protected

  def user_response(user)
    {
      id: user.id,
      email: user.email,
      name: user.name,
      role: user.role,
      created_at: user.created_at
    }
  end

  # Allow role to be passed only if current_user is admin
  def sign_up_params
    if current_user&.admin?
      params.require(:user).permit(:email, :password, :password_confirmation, :name, :role)
    else
      params.require(:user).permit(:email, :password, :password_confirmation, :name)
    end
  end


  def account_update_params
    params.require(:user).permit(:email, :password, :password_confirmation, :name)
  end
end
