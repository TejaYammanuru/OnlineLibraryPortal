class Users::SessionsController < Devise::SessionsController
  include ActionController::MimeResponds
  respond_to :json

  
  def create
    begin
      self.resource = warden.authenticate!(auth_options)
      sign_in(resource_name, resource)

      log_auth_event(resource, 'login')

      render json: {
        message: 'Logged in successfully',
        user: user_response(resource)
      }, status: :ok
    rescue => e
      render json: { message: 'Invalid email or password' }, status: :unauthorized
    end
  end

  
  def destroy
    user = current_user
    sign_out(resource_name)

    log_auth_event(user, 'logout') if user.present?

    render json: { message: 'Logged out successfully' }, status: :ok
  end

  
  def verify_token
    token = request.headers['Authorization']&.split(' ')&.last

    if token.blank?
      render json: { error: 'Token missing' }, status: :unauthorized and return
    end

    begin
      secret_key = 'Teja'
      decoded_token = JWT.decode(token, secret_key, true, algorithm: 'HS256')
      payload = decoded_token.first

      user_id = payload['sub'] || payload['id']
      user = User.find_by(id: user_id)

      if user
        render json: {
          id: user.id,
          email: user.email,
          name: user.name,
          role: user.role
        }, status: :ok
      else
        render json: { error: 'User not found' }, status: :unauthorized
      end

    rescue JWT::DecodeError => e
      render json: { error: "Invalid token: #{e.message}" }, status: :unauthorized
    end
  end

  
  def profile
    token = request.headers['Authorization']&.split(' ')&.last
    if token.blank?
      render json: { error: 'Token missing' }, status: :unauthorized and return
    end

    begin
      secret_key = 'Teja'
      decoded_token = JWT.decode(token, secret_key, true, algorithm: 'HS256')
      payload = decoded_token.first

      user_id = payload['sub'] || payload['id']
      user = User.find_by(id: user_id)

      if user
        render json: {
          user: {
            id: user.id,
            name: user.name,
            email: user.email,
            role: user.role
          }
        }, status: :ok
      else
        render json: { error: 'User not found' }, status: :unauthorized
      end

    rescue JWT::DecodeError => e
      render json: { error: "Invalid token: #{e.message}" }, status: :unauthorized
    end
  end

  private

  def user_response(user)
    {
      id: user.id,
      email: user.email,
      name: user.name,
      role: user.role
    }
  end

  def log_auth_event(user, action)
    return unless user.present?

    AuthLog.create!(
      user_id: user.id,
      name: user.name,
      email: user.email,
      role: user.role,
      action: action
    )
  end
end
