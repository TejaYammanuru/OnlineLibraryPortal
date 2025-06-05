# frozen_string_literal: true

class Users::RegistrationsController < Devise::RegistrationsController
  respond_to :json

  # POST /signup
  def create
    build_resource(sign_up_params)

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
      Rails.logger.error("Signup failed: #{resource.errors.full_messages}")
      render json: {
        message: 'User registration failed',
        errors: resource.errors.full_messages
      }, status: :unprocessable_entity
    end
  end

  
  def update
    self.resource = resource_class.to_adapter.get!(send(:"current_#{resource_name}").to_key)
    
    if resource.update(account_update_params)
      render json: {
        message: 'User updated successfully',
        user: user_response(resource)
      }, status: :ok
    else
      Rails.logger.error("Update failed: #{resource.errors.full_messages}")
      render json: {
        message: 'User update failed',
        errors: resource.errors.full_messages
      }, status: :unprocessable_entity
    end
  end

  
  def destroy
    resource = current_user
    resource.destroy
    render json: { message: 'Account deleted successfully' }, status: :ok
  end

  protected
  def user_response(user)
    {
      id: user.id,
      email: user.email,
      name: user.name,
      created_at: user.created_at
    }
  end
end
