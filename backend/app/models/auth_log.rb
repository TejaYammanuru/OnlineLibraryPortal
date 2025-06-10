
class AuthLog
  include Mongoid::Document
  include Mongoid::Timestamps

  field :user_id, type: Integer
  field :name, type: String
  field :email, type: String
  field :role, type: String
  field :action, type: String  
end
