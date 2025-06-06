class Users::SessionsController < Devise::SessionsController
  include ActionController::MimeResponds 
  respond_to :json

 
  def create
    self.resource = warden.authenticate!(auth_options)

    if resource
      sign_in(resource_name, resource)
      render json: {
        message: 'Logged in successfully',
        user: user_response(resource)
      }, status: :ok
    else
      render json: {
        message: 'Invalid email or password'
      }, status: :unauthorized
    end
  end

 
  def destroy
    sign_out(resource_name)
    render json: { message: 'Logged out successfully' }, status: :ok
  end

  private

  def user_response(user)
    {
      id: user.id,
      email: user.email,
      name: user.name
    }
  end
end
